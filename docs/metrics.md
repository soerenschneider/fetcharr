# Metrics

All metrics are namespaced with `fetcharr`.

## Subsystem Stats
| Metric Name                         | Description                               | Type     |
|-------------------------------------|-------------------------------------------|----------|
| stats_num_files                     | Number of files in rsync stats            | Gauge    |
| stats_num_created_files             | Number of created files in rsync stats    | Gauge    |
| stats_num_deleted_files             | Number of deleted files in rsync stats    | Gauge    |
| stats_num_regular_files_transferred | Number of regular files transferred       | Gauge    |
| stats_total_file_size_bytes         | Total file size in bytes in rsync stats   | Gauge    |
| stats_total_bytes_sent              | Total bytes sent in rsync stats           | Gauge    |
| stats_total_bytes_received          | Total bytes received in rsync stats       | Gauge    |

## Subsystem Runner
| Metric Name                    | Description                                   | Type       | Labels |
|--------------------------------|-----------------------------------------------|------------|--------|
| runner_sync_seconds            | Synchronization time in seconds               | Gauge      |        |
| runner_hook_error_total        | Hook errors total count (by type)             | GaugeVec   | type   |
| runner_sync_timestamp_seconds  | Most recent synchronization timestamp         | Gauge      |        |
| runner_sync_errors_total       | Runner errors total count (by implementation) | GaugeVec   | impl   |

## Subsystem Events
| Metric Name                          | Description                               | Type       | Labels  |
|--------------------------------------|-------------------------------------------|------------|---------|
| events_received_timestamp_seconds    | Event received timestamp in seconds       | GaugeVec   | source  |
| events_received_total                | Total events received count (by source)   | GaugeVec   | source  |
