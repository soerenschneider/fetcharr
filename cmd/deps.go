package main

import (
	"errors"
	"time"

	"github.com/soerenschneider/fetcharr/internal"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/events"
	"github.com/soerenschneider/fetcharr/internal/syncer"
	"github.com/soerenschneider/fetcharr/internal/syncer/rsync"
)

func buildSyncer(conf *config.Config) (syncer.Syncer, error) {
	if conf.SyncerImpl == "rsync" {
		return buildRsync(conf)
	}

	return nil, errors.New("unknown syncer")
}

func buildApp(conf *config.Config, syncImpl syncer.Syncer, events chan events.EventSyncRequest) (*internal.Fetcharr, error) {
	var opts []internal.FetcharrOpts
	opts = append(opts, internal.WithTimeout(6*time.Hour))
	if len(conf.Hooks) > 0 {
		hooks, err := buildHooks(conf.Hooks)
		if err != nil {
			return nil, err
		}
		opts = append(opts, internal.WithHooks(hooks))
	}

	return internal.NewFetcharr(syncImpl, events, opts...)
}

func buildRsync(conf *config.Config) (*rsync.Rsync, error) {
	var opts []rsync.RsyncOpt

	if len(conf.Rsync.BandwidthLimit) > 0 {
		opts = append(opts, rsync.BandwidthLimit(conf.Rsync.BandwidthLimit))
	}

	if len(conf.Rsync.ExcludePattern) > 0 {
		opts = append(opts, rsync.Exclude(conf.Rsync.ExcludePattern))
	}

	if conf.Rsync.RemoveSourceFiles {
		opts = append(opts, rsync.RemoveSourceFiles())
	}

	if len(conf.Rsync.RemoteShell) > 0 {
		opts = append(opts, rsync.RemoteShell(conf.Rsync.RemoteShell))
	}

	return rsync.New(conf.Rsync.Host, conf.Rsync.RemoteDir, conf.Rsync.LocalDir, opts...)
}
