package events

import (
	"context"
	"sync"
)

type EventSource interface {
	Listen(ctx context.Context, events chan EventSyncRequest, wg *sync.WaitGroup) error
}

type EventSyncRequest struct {
	Source   string
	Metadata string
	Response chan error
}
