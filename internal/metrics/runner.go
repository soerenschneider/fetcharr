package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const runnerSubsystem = "runner"

var (
	SyncTime = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: runnerSubsystem,
		Name:      "sync_seconds",
	})

	HookRuns = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: runnerSubsystem,
		Name:      "hook_runs_total",
	}, []string{"type", "name"})

	HookError = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: runnerSubsystem,
		Name:      "hook_error_total",
	}, []string{"type", "name"})

	MostRecentSync = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: runnerSubsystem,
		Name:      "sync_timestamp_seconds",
	})

	RunnerErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: runnerSubsystem,
		Name:      "sync_errors_total",
	}, []string{"impl"})
)
