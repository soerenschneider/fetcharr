package ticker

import (
	"context"
	"errors"
	"sync"
	"time"
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

func (t *Ticker) Listen(ctx context.Context, events chan bool, wg *sync.WaitGroup) error {
	wg.Add(1)
	ticker := time.NewTicker(t.interval)

	for {
		select {
		case <-ticker.C:
			events <- true
		case <-ctx.Done():
			ticker.Stop()
			wg.Done()
			return nil
		}
	}
}
