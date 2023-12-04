package runner

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/hooks"
	"github.com/soerenschneider/fetcharr/internal/metrics"
	"github.com/soerenschneider/fetcharr/internal/syncer"
	"go.uber.org/multierr"
)

const defaultSyncerTimeout = 12 * time.Hour

type Runner struct {
	syncer        syncer.Syncer
	wantsSync     atomic.Bool
	once          sync.Once
	syncerTimeout time.Duration

	hooks hooks.HookExectuor
}

type RunnerOpts func(*Runner) error

func New(syncer syncer.Syncer, opts ...RunnerOpts) (*Runner, error) {
	if syncer == nil {
		return nil, errors.New("emtpy syncer supplied")
	}

	runner := &Runner{
		syncer:        syncer,
		syncerTimeout: defaultSyncerTimeout,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(runner); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return runner, errs
}

func (r *Runner) QueueSync() {
	swapped := r.wantsSync.CompareAndSwap(false, true)
	if swapped {
		log.Debug().Msg("Setting wantsSync=true")
	}
}

func (r *Runner) Start(ctx context.Context, wg *sync.WaitGroup) error {
	if ctx == nil {
		return errors.New("empty context supplied")
	}

	if wg == nil {
		return errors.New("empty waitgroup supplied")
	}

	r.once.Do(func() {
		wg.Add(1)
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			r.work(ctx)
			for {
				select {
				case <-ticker.C:
					r.work(ctx)
				case <-ctx.Done():
					log.Info().Msg("Stopping runner...")
					ticker.Stop()
					wg.Done()
					return
				}
			}
		}()
	})

	return nil
}

func (r *Runner) work(ctx context.Context) {
	log.Debug().Msg("Checking if a new run is needed")
	if !r.wantsSync.CompareAndSwap(true, false) {
		log.Debug().Msg("No request to sync")
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
