package internal

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/events"
	"github.com/soerenschneider/fetcharr/internal/hooks"
	"github.com/soerenschneider/fetcharr/internal/metrics"
	"github.com/soerenschneider/fetcharr/internal/syncer"
	"go.uber.org/multierr"
)

const defaultSyncerTimeout = 12 * time.Hour
const defaultCooldownDuration = 30 * time.Second

type Fetcharr struct {
	syncer        syncer.Syncer
	wantsSync     atomic.Bool
	once          sync.Once
	syncerTimeout time.Duration

	cooldownTimer time.Duration

	mut       sync.Mutex
	eventChan chan events.EventSyncRequest

	hooks hooks.HookExectuor
}

type FetcharrOpts func(*Fetcharr) error

func NewFetcharr(syncer syncer.Syncer, events chan events.EventSyncRequest, opts ...FetcharrOpts) (*Fetcharr, error) {
	if syncer == nil {
		return nil, errors.New("emtpy syncer supplied")
	}

	runner := &Fetcharr{
		syncer:        syncer,
		syncerTimeout: defaultSyncerTimeout,
		eventChan:     events,
		cooldownTimer: defaultCooldownDuration,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(runner); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return runner, errs
}

func (r *Fetcharr) Start(ctx context.Context, wg *sync.WaitGroup) error {
	if wg == nil {
		return errors.New("empty waitgroup supplied")
	}

	r.once.Do(func() {
		defer wg.Done()
		wg.Add(1)
		r.work(ctx)

		for {
			select {
			case event := <-r.eventChan:
				metrics.EventsReceived.WithLabelValues(event.Source).Inc()
				metrics.EventReceivedTimestamp.WithLabelValues(event.Source).SetToCurrentTime()

				if r.wantsSync.CompareAndSwap(false, true) {
					select {
					case event.Response <- nil:
					case <-time.After(1 * time.Second):
						log.Warn().Str("event_source", event.Source).Msg("hanging goroutine / empty channel")
					}

					log.Debug().Msg("Setting wantsSync=true")
					go func() {
						r.work(ctx)
					}()
				} else {
					metrics.EventsIgnored.WithLabelValues(event.Source).Inc()

					select {
					case event.Response <- errors.New("cooldown phase"):
					case <-time.After(5 * time.Second):
						log.Error().Str("event_source", event.Source).Msg("hanging goroutine / empty channel")
					}
				}
			case <-ctx.Done():
				log.Info().Msg("Stopping runner...")
				return
			}
		}
	})

	return nil
}

func (r *Fetcharr) work(ctx context.Context) {
	r.mut.Lock()
	defer r.mut.Unlock()

	log.Info().Msg("Checking if a new run is needed")
	if r.wantsSync.Load() {
		log.Info().Msg("Syncing and entering cooldown phase")

		cooldownCtx, cancel := context.WithTimeout(ctx, r.cooldownTimer)
		go func() {
			defer cancel()
			<-cooldownCtx.Done()
			r.wantsSync.Store(false)
		}()
	} else {
		log.Info().Msg("No request to sync")
		return
	}

	metrics.MostRecentSync.SetToCurrentTime()
	ctx, cancel := context.WithTimeout(ctx, r.syncerTimeout)
	defer cancel()

	if r.hooks.HasHooksDefined() {
		err := r.hooks.Run(ctx, config.PRE, nil)
		if err != nil {
			log.Error().Err(err).Msg("errors while running pre hooks")
		}
	}

	log.Info().Msg("Starting sync")
	timeStarted := time.Now()
	stats, err := r.syncer.Sync(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error while syncing")
		metrics.RunnerErrors.WithLabelValues(r.syncer.Name()).Inc()
	} else {
		timeTotal := time.Since(timeStarted)
		metrics.SyncTime.Set(timeTotal.Seconds())
		log.Info().Msgf("Finished sync, transferred %d files (%s) in %v", stats.NumTransferredFiles, stats.TotalBytesReceivedHumanized, timeTotal)
	}

	if r.hooks.HasHooksDefined() {
		err := r.hooks.Run(ctx, detectStage(err, stats), stats)
		if err != nil {
			log.Error().Err(err).Msg("errors while running post hooks")
		}
	}
}

func detectStage(err error, stats *syncer.Stats) config.Stage {
	if err != nil {
		return config.POST_FAILURE
	}

	if stats.NumTransferredFiles > 0 {
		return config.POST_SUCCESS_TRANSFER
	}

	return config.POST_SUCCESS
}
