package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soerenschneider/fetcharr/internal/syncer"
)

const statsSubsystem = "stats"

var (
	numFilesMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: statsSubsystem,
		Name:      "num_files",
		Help:      "Number of files in rsync stats",
	})
	numCreatedFilesMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: statsSubsystem,
		Name:      "num_created_files",
		Help:      "Number of created files in rsync stats",
	})
	numDeletedFilesMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: statsSubsystem,
		Name:      "num_deleted_files",
		Help:      "Number of deleted files in rsync stats",
	})
	numRegularFilesTransferredMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: statsSubsystem,
		Name:      "num_regular_files_transferred",
		Help:      "Number of regular files transferred in rsync stats",
	})
	totalFileSizeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: statsSubsystem,
		Name:      "total_file_size_bytes",
		Help:      "Total file size in bytes in rsync stats",
	})
	totalBytesSentMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: statsSubsystem,
		Name:      "total_bytes_sent",
		Help:      "Total bytes sent in rsync stats",
	})
	totalBytesReceivedMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: statsSubsystem,
		Name:      "total_bytes_received",
		Help:      "Total bytes received in rsync stats",
	})
)

func UpdatePrometheusMetrics(stats syncer.Stats) {
	numFilesMetric.Set(float64(stats.NumFiles))
	numCreatedFilesMetric.Set(float64(stats.NumCreatedFiles))
	numDeletedFilesMetric.Set(float64(stats.NumDeletedFiles))
	numRegularFilesTransferredMetric.Set(float64(stats.NumTransferredFiles))
	totalFileSizeMetric.Set(float64(stats.TotalFileSize))
	totalBytesSentMetric.Set(float64(stats.TotalBytesSent))
	totalBytesReceivedMetric.Set(float64(stats.TotalBytesReceived))
}
