package syncer

import (
	"context"
)

type Stats struct {
	Files               []string
	NumFiles            int64
	NumCreatedFiles     int64
	NumDeletedFiles     int64
	NumTransferredFiles int64
	TotalFileSize       int64
	TotalBytesSent      int64
	TotalBytesReceived  int64

	TotalFileSizeHumanized      string
	TotalBytesSentHumanized     string
	TotalBytesReceivedHumanized string
}

type Syncer interface {
	NeedsSync(ctx context.Context) bool
	Sync(ctx context.Context) (*Stats, error)
	Name() string
}
