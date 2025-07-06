package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/events"
	"github.com/soerenschneider/fetcharr/internal/metrics"
	"github.com/soerenschneider/fetcharr/internal/syncer"
)

var (
	configFile   string
	printVersion bool
	debug        bool
)

type Deps struct {
	conf config.Config

	eventsChan chan events.EventSyncRequest
	wg         *sync.WaitGroup

	app *internal.Fetcharr

	// movable parts
	syncImpl     syncer.Syncer
	eventSources []events.EventSource
}

const (
	defaultConfigFile = "/etc/fetcharr.yaml"
)

func main() {
	// parse flags
	parseFlags()
	if printVersion {
		//nolint: forbidigo
		fmt.Println(internal.BuildVersion)
		os.Exit(0)
	}

	log.Info().Msgf("Starting fetcharr, version %s (%s)", internal.BuildVersion, internal.CommitHash)

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// parse config
	conf, err := config.Read(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}

	if err := config.Validate(conf); err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}

	// build deps
	deps := &Deps{}
	deps.conf = *conf
	deps.eventsChan = make(chan events.EventSyncRequest, 12)
	deps.wg = &sync.WaitGroup{}

	deps.syncImpl, err = buildSyncer(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build syncer")
	}

	deps.app, err = buildApp(conf, deps.syncImpl, deps.eventsChan)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build runner")
	}

	deps.eventSources, err = buildEventSources(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build event sources")
	}

	run(deps)
}

func run(deps *Deps) {
	ctx, cancel := context.WithCancel(context.Background())

	for _, eventSource := range deps.eventSources {
		go func(source events.EventSource) {
			err := source.Listen(ctx, deps.eventsChan, deps.wg)
			if err != nil {
				log.Error().Err(err).Msg("listening on event source failed")
			}
		}(eventSource)
	}

	go func() {
		err := deps.app.Start(ctx, deps.wg)
		if err != nil {
			log.Fatal().Err(err).Msg("could not start runner")
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Info().Msg("Caught signal, shutting down gracefully")
		cancel()
	}()

	if len(deps.conf.MetricsAddr) > 0 {
		go metrics.StartServer(deps.conf.MetricsAddr)
	}

	if err := deps.app.Start(ctx, deps.wg); err != nil {
		log.Error().Err(err).Msg("error running loop")
	}

	deps.wg.Wait()
	log.Info().Msg("All components shut down, bye!")
}

func parseFlags() {
	flag.StringVar(&configFile, "config", defaultConfigFile, fmt.Sprintf("Path to the config file (default %s)", defaultConfigFile))
	flag.BoolVar(&printVersion, "version", false, "Print printVersion and exit")
	flag.BoolVar(&debug, "debug", false, "Print debug information")
	flag.Parse()
}
