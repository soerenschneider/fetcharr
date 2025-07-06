package internal

import (
	"errors"
	"time"

	"github.com/soerenschneider/fetcharr/internal/hooks"
)

func WithTimeout(timeout time.Duration) func(*Fetcharr) error {
	return func(r *Fetcharr) error {
		r.syncerTimeout = timeout
		return nil
	}
}

func WithHooks(h []hooks.Hook) FetcharrOpts {
	return func(r *Fetcharr) error {
		if len(h) == 0 {
			return errors.New("empty pre hooks supplied")
		}

		hookExecutor, err := hooks.NewHookExecutor(h...)
		if err != nil {
			return err
		}

		r.hooks = *hookExecutor
		return nil
	}
}
