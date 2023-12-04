package metrics

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog/log"
)

const namespace = "fetcharr"

func StartServer(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msgf("Can not start metrics server")
	}
}

func StartPusher(url string) {
	ticker := time.NewTicker(60 * time.Second)

	for range ticker.C {
		err := Push(url)
		if err != nil {
			log.Error().Err(err).Msg("could not push metrics to remote ")
		}
	}
}

func Push(url string) error {
	err := push.New(url, "fetcharr").
		Gatherer(prometheus.DefaultGatherer).
		Push()

	return err
}
