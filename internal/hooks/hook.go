package hooks

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/metrics"
	"github.com/soerenschneider/fetcharr/internal/syncer"
	"go.uber.org/multierr"
)

type Hook interface {
	Run(ctx context.Context, stats *syncer.Stats) error
	config.HookConfig
}

type HookExectuor struct {
	hooks []Hook
}

func NewHookExecutor(hooks ...Hook) (*HookExectuor, error) {
	return &HookExectuor{
		hooks: hooks,
	}, nil
}

func (h *HookExectuor) HasHooksDefined() bool {
	return len(h.hooks) > 0
}

func (h *HookExectuor) Run(ctx context.Context, stage config.Stage, stats *syncer.Stats) error {
	var errs error
	for _, hook := range h.hooks {
		if hook.GetStage() == stage {
			metrics.HookRuns.WithLabelValues(string(stage), hook.GetName()).Inc()

			log.Info().Str("hook", hook.GetName()).Msg("Running hook")
			if err := hook.Run(ctx, stats); err != nil {
				metrics.HookError.WithLabelValues(string(stage), hook.GetName()).Inc()
				errs = multierr.Append(errs, err)
				if hook.ExitOnErr() && stage == config.PRE {
					log.Fatal().Err(err).Str("hook", hook.GetName()).Msg("exiting due to failed pre-hook marked as 'exit on error'")
				}
			}
		}
	}

	return errs
}
