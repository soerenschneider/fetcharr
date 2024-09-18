package webhook_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

const defaultPath = "/fetcharr-webhook_server"

type WebhookServer struct {
	address string

	// optional
	path     string
	certFile string
	keyFile  string
}

type WebhookServerOpts func(*WebhookServer) error

func New(address string, opts ...WebhookServerOpts) (*WebhookServer, error) {
	if len(address) == 0 {
		return nil, errors.New("empty address provided")
	}

	w := &WebhookServer{
		address: address,
		path:    defaultPath,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(w); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return w, errs
}

func (w *WebhookServer) IsTLSConfigured() bool {
	return len(w.certFile) > 0 && len(w.keyFile) > 0
}

func (w *WebhookServer) Listen(ctx context.Context, events chan bool, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()
	mux := http.NewServeMux()

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events <- true
		w.WriteHeader(http.StatusOK)
	}
	mux.HandleFunc(w.path, handler)

	server := http.Server{
		Addr:              w.address,
		Handler:           mux,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	errChan := make(chan error)
	go func() {
		if w.IsTLSConfigured() {
			if err := server.ListenAndServeTLS(w.certFile, w.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- fmt.Errorf("can not start webhook_server server: %w", err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- fmt.Errorf("can not start webhook_server server: %w", err)
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("Stopping webhook_server server")
		err := server.Shutdown(ctx)
		return err
	case err := <-errChan:
		return err
	}
}
