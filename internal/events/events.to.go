package events

import (
	"context"
	"sync"
)

type EventSource interface {
	Listen(ctx context.Context, events chan bool, wg *sync.WaitGroup) error
}
