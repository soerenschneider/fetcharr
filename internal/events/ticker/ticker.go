package ticker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/events"
)

type Ticker struct {
	interval time.Duration
}

func NewTicker(interval time.Duration) (*Ticker, error) {
	if interval.Seconds() <= 300 {
		return nil, errors.New("interval must not be <= 300 seconds")
	}

	return &Ticker{interval: interval}, nil
}

func (t *Ticker) Listen(ctx context.Context, eventChan chan events.EventSyncRequest, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()
	ticker := time.NewTicker(t.interval)

	work := func() {
		req := events.EventSyncRequest{
			Source:   "ticker",
			Metadata: "",
			Response: make(chan error),
		}
		eventChan <- req

		select {
		case err := <-req.Response:
			if err != nil {
				log.Error().Str("component", "ticker").Err(err).Msg("received error response from fetcharr")
			}
		case <-time.After(3 * time.Second):
			log.Warn().Str("component", "ticker").Msgf("timeout waiting for goroutine")
		}
	}

	work()

	for {
		select {
		case <-ticker.C:
			work()
		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}
