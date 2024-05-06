package hooks

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/syncer"
	"go.uber.org/multierr"
)

type CmdHook struct {
	conf config.CmdHookConfig
}

func NewCmdHook(conf config.CmdHookConfig) (*CmdHook, error) {
	return &CmdHook{conf: conf}, nil
}

func (c *CmdHook) Run(ctx context.Context, _ *syncer.Stats) error {
	var errs error
	for index, hook := range c.conf.Cmds {
		if !*c.conf.StopOnError || errs == nil {
			if len(c.conf.Cmds) > 1 {
				log.Info().Msgf("Running cmd #%d: %s", index, hook)
			}

			split := strings.Split(hook, " ")
			var args []string
			if len(split) > 1 {
				args = split[1:]
			}

			cmd := exec.CommandContext(ctx, split[0], args...) // #nosec: G204
			var stdErr bytes.Buffer
			cmd.Stderr = &stdErr

			err := cmd.Run()
			if err != nil {
				errs = multierr.Append(errs, fmt.Errorf("hook %q, cmd %q (#%d): %s", c.conf.Name, args[0], index, stdErr.String()))
			}
		}
	}

	return errs
}

func (w *CmdHook) GetType() string {
	return w.conf.GetType()
}

func (w *CmdHook) GetName() string {
	return w.conf.GetName()
}

func (w *CmdHook) GetStage() config.Stage {
	return w.conf.GetStage()
}

func (c *CmdHook) ExitOnErr() bool {
	return c.conf.ExitOnErr()
}
