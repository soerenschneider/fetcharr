package main

import (
	"net/http"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/hooks"
	"go.uber.org/multierr"
)

var (
	once   sync.Once
	client *http.Client
)

func buildHooks(hooksConf []config.HookConfigContainer) ([]hooks.Hook, error) {
	var errs error
	var builtHooks []hooks.Hook
	for _, h := range hooksConf {
		var hook hooks.Hook
		var err error

		switch h.HookConfig.GetType() {
		case config.WebhookType:
			conf := h.HookConfig.(*config.WebhookConfig)
			hook, err = buildWebhookClient(*conf)
			log.Debug().Msgf("Built webhook %q for stage %q", hook.GetName(), hook.GetStage())
		case config.CmdHookType:
			conf := h.HookConfig.(*config.CmdHookConfig)
			hook, err = hooks.NewCmdHook(*conf)
			log.Debug().Msgf("Built cmd %q for stage %q", hook.GetName(), hook.GetStage())
		}

		if err != nil {
			errs = multierr.Append(errs, err)
		} else {
			builtHooks = append(builtHooks, hook)
		}
	}

	return builtHooks, errs
}

func buildWebhookClient(conf config.WebhookConfig) (*hooks.WebhookClient, error) {
	once.Do(func() {
		c := retryablehttp.NewClient()
		c.Logger = &ZerologAdapter{}
		client = c.StandardClient()
	})

	return hooks.NewWebhookClient(client, conf)
}

type ZerologAdapter struct {
}

// Debug logs a debug-level message
func (z *ZerologAdapter) Debug(msg string, keysAndValues ...interface{}) {
	log.Debug().Str("checker", "prometheus").Interface("details", keysAndValues).Msg(msg)
}

// Info logs an info-level message
func (z *ZerologAdapter) Info(msg string, keysAndValues ...interface{}) {
	log.Info().Str("checker", "prometheus").Interface("details", keysAndValues).Msg(msg)
}

// Warn logs a warning-level message
func (z *ZerologAdapter) Warn(msg string, keysAndValues ...interface{}) {
	log.Warn().Str("checker", "prometheus").Interface("details", keysAndValues).Msg(msg)
}

// Error logs an error-level message
func (z *ZerologAdapter) Error(msg string, keysAndValues ...interface{}) {
	log.Error().Str("checker", "prometheus").Interface("details", keysAndValues).Msg(msg)
}
