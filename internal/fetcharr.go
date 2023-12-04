package internal

import (
	"context"
	"errors"
	"sync"

	"github.com/soerenschneider/fetcharr/internal/metrics"
	"github.com/soerenschneider/fetcharr/internal/runner"

	"github.com/rs/zerolog/log"
)

type Fetcharr struct {
	eventChan chan bool
	runner    *runner.Runner
}

func NewFetcharr(eventChan chan bool, runner *runner.Runner) (*Fetcharr, error) {
	if eventChan == nil {
		return nil, errors.New("closed channel supplied")
	}

	if runner == nil {
		return nil, errors.New("empty syncer supplied")
	}

	return &Fetcharr{
		eventChan: eventChan,
		runner:    runner,
	}, nil
}

func (a *Fetcharr) Loop(ctx context.Context, wg *sync.WaitGroup) error {
	if wg == nil {
		return errors.New("empty waitgroup supplied")
	}
	wg.Add(1)

	a.runner.QueueSync()

	for {
		select {
		case <-a.eventChan:
			// TODO: set source
			metrics.EventsReceived.WithLabelValues("unknown").Inc()
			metrics.EventReceivedTimestamp.WithLabelValues("unknown").SetToCurrentTime()
			log.Info().Msg("Event is queueing new sync")
			a.runner.QueueSync()
		case <-ctx.Done():
			log.Info().Msg("Stopping fetcharr...")
			wg.Done()
			return nil
		}
	}
}
