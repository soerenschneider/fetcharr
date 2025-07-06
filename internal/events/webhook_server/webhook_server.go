package webhook_server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/events"
	"go.uber.org/multierr"
)

const defaultPath = "/fetcharr-webhook_server"

type WebhookServer struct {
	address   string
	eventChan chan events.EventSyncRequest

	shutdownOngoing atomic.Bool

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

func (s *WebhookServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	syncRequest := events.EventSyncRequest{
		Source:   "webhook",
		Metadata: getIP(r),
		Response: make(chan error),
	}

	if s.shutdownOngoing.Load() {
		http.Error(w, "Server is shutting down", http.StatusServiceUnavailable)
		return
	}

	select {
	case s.eventChan <- syncRequest:
		// proceed
	case <-time.After(500 * time.Millisecond):
		http.Error(w, "Server busy, try again later", http.StatusServiceUnavailable)
		return
	}

	select {
	case err := <-syncRequest.Response:
		if err != nil {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case <-time.After(5 * time.Second):
		http.Error(w, "Timeout waiting for event processing", http.StatusGatewayTimeout)
		return
	}
}

func (w *WebhookServer) Listen(ctx context.Context, eventChan chan events.EventSyncRequest, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	w.eventChan = eventChan

	mux := http.NewServeMux()
	mux.HandleFunc(w.path, w.handler)

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
		w.shutdownOngoing.Store(true)

		log.Info().Msg("Stopping webhook_server server")
		err := server.Shutdown(ctx)
		return err
	case err := <-errChan:
		return err
	}
}

func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xrip := r.Header.Get("X-Real-IP")
	if xrip != "" {
		return xrip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
