package syncer

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

type DummySyncer struct {
}

func (r *DummySyncer) NeedsSync(ctx context.Context) bool {
	return true
}

func (s *DummySyncer) Sync(ctx context.Context) error {
	log.Info().Msg("Sync is running")
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	sleep := r.Intn(15)
	log.Info().Msgf("Going to sleep for %d seconds", sleep)
	cmd := exec.CommandContext(ctx, "sleep", fmt.Sprintf("%d", sleep)) // #nosec G204
	err := cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error running syncer")
		return err
	}
	log.Info().Msg("Sync is finished")
	return nil
}
