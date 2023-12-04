package rsync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/metrics"
	"github.com/soerenschneider/fetcharr/internal/syncer"
	"github.com/soerenschneider/fetcharr/pkg"
	"go.uber.org/multierr"
)

const rsyncOutFormatPrefix = "downloaded="

type Rsync struct {
	cmd  string
	args []string
}

type RsyncOpt func(*Rsync) error

func New(host string, remoteDir string, localDir string, opts ...RsyncOpt) (*Rsync, error) {
	if len(host) == 0 {
		return nil, errors.New("rsync: empty host provided")
	}
	if len(remoteDir) == 0 {
		return nil, errors.New("rsync: empty remoteDir provided")
	}
	if len(localDir) == 0 {
		return nil, errors.New("rsync: empty localDir provided")
	}

	args := []string{
		"--partial",
		"--recursive",
		"--stats",
		"--out-format=" + rsyncOutFormatPrefix + "%n",
		fmt.Sprintf("%s:%s", host, remoteDir),
		localDir,
	}

	rsync := &Rsync{
		cmd:  "rsync",
		args: args,
	}

	var errs error
	for _, opt := range opts {
		err := opt(rsync)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return rsync, errs
}

func (r *Rsync) Name() string {
	return "rsync"
}

func (r *Rsync) String() string {
	return fmt.Sprintf("%s %s", r.cmd, r.args)
}

func (r *Rsync) NeedsSync(ctx context.Context) bool {
	return true
}

func (r *Rsync) Sync(ctx context.Context) (*syncer.Stats, error) {
	log.Info().Msgf("Running rsync using args %v", r.args)
	cmd := exec.CommandContext(ctx, r.cmd, r.args...) // #nosec G204

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("running rsync failed: %s", stderr.String())
	}

	parsed := parseRsyncOutput(stdout.String())
	return &parsed, nil
}

// nolint: cyclop
func parseRsyncOutput(output string) syncer.Stats {
	var stats syncer.Stats

	lines := strings.Split(strings.ToLower(output), "\n")
	var files []string
	for _, line := range lines {
		if strings.HasPrefix(line, rsyncOutFormatPrefix) {
			file := line[len(rsyncOutFormatPrefix):]
			file = filepath.Base(file)
			files = append(files, file)
		} else {
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])

			switch key {
			case "number of files":
				stats.NumFiles = extractInt(valueStr)
			case "number of files transferred": // mac osx
				stats.NumTransferredFiles = extractInt(valueStr)
			case "number of regular files transferred": // linux
				stats.NumTransferredFiles = extractInt(valueStr)
			case "number of created files":
				stats.NumCreatedFiles = extractInt(valueStr)
			case "number of deleted files":
				stats.NumDeletedFiles = extractInt(valueStr)
			case "total file size":
				stats.TotalFileSize = extractInt(valueStr)
			case "total bytes sent":
				stats.TotalBytesSent = extractInt(valueStr)
			case "total bytes received":
				stats.TotalBytesReceived = extractInt(valueStr)
			}
		}
	}

	stats.Files = files
	metrics.UpdatePrometheusMetrics(stats)
	stats.TotalFileSizeHumanized = pkg.BytesToHumanSize(stats.TotalFileSize)
	stats.TotalBytesSentHumanized = pkg.BytesToHumanSize(stats.TotalBytesSent)
	stats.TotalBytesReceivedHumanized = pkg.BytesToHumanSize(stats.TotalBytesReceived)
	return stats
}

func extractInt(valueStr string) int64 {
	// rsync on linux returned values with a ',' separator
	floatVal := extractFloat(valueStr)
	return int64(floatVal)
}

func extractFloat(valueStr string) float64 {
	valueStr = strings.Replace(valueStr, ",", "", -1)
	value, err := strconv.ParseFloat(strings.Split(valueStr, " ")[0], 64)
	if err != nil {
		return 0.0
	}
	return value
}
