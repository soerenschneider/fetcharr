# Configuration

## Config Reference
| Field Name                | Description                                                                                 | YAML Key                | Data Type            | Optional |
|---------------------------|---------------------------------------------------------------------------------------------|-------------------------|----------------------|----------|
| `SyncerImpl`              | The syncer implementation (must be "rsync").                                                | `"syncer_impl"`         | String               | No       |
| `PreHooks`                | An array of pre-sync hooks (optional).                                                      | `"pre_hooks"`           | Array of strings     | Yes      |
| `PostHooks`               | An array of post-sync hooks (optional).                                                     | `"post_hooks"`          | Array of strings     | Yes      |
| `EventSourceImpl`         | An array of event source implementations. Valid choices are "kafka", "webhook" and "ticker" | `"events_impl"`         | Array of strings     | Yes      |
| `MetricsAddr`             | The metrics server address (optional, must be in the format "host:port").                   | `"metrics_addr"`        | String (TCP Address) | Yes      |

### Rsync Object
| Field Name                | Description                                                          | YAML Key                | Data Type            | Optional    |
|---------------------------|----------------------------------------------------------------------|-------------------------|----------------------|-------------|
| `Rsync`                   | Configuration for the Rsync syncer.                                  | `"rsync"`               | Object               | No          |
| `Rsync.Host`              | The Rsync host.                                                      | `"host"`                | String               | Yes         |
| `Rsync.LocalDir`          | The local directory (must end with "/").                             | `"local_dir"`           | String               | No          |
| `Rsync.RemoteDir`         | The remote directory (must end with "/").                            | `"remote_dir"`          | String               | No          |
| `Rsync.BandwidthLimit`    | The bandwidth limit for Rsync.                                       | `"bwlimit"`             | String               | Yes         |
| `Rsync.ExcludePattern`    | The pattern for excluding files during synchronization.              | `"exclude"`             | String               | Yes         |
| `Rsync.RemoveSourceFiles` | Flag to remove source files after successful sync.                   | `"remove_source_files"` | Boolean              | Yes         |
| `Rsync.RemoteShell`       | The remote shell command (optional). Can be used to set ssh options. | `"remote_shell"`        | String               | Yes         |

### Kafka Object
| Field Name                | Description                                                                   | YAML Key                | Data Type            | Optional    |
|---------------------------|-------------------------------------------------------------------------------|-------------------------|----------------------|-------------|
| `Kafka`                   | Configuration for `Kafka event source.                                        | `"kafka"`               | Object               | Yes         |
| `Kafka.Brokers`           | The list of Kafka brokers (required if `"kafka"` is in `EventSourceImpl`).    | `"brokers"`             | Array of strings     | Yes         |
| `Kafka.Topic`             | The Kafka topic (required if `"kafka"` is in `EventSourceImpl`).              | `"topic"`               | String               | Yes         |
| `Kafka.GroupId`           | The Kafka consumer group ID (required if `"kafka"` is in `EventSourceImpl`).  | `"group_id"`            | String               | Yes         |
| `Kafka.Partition`         | The Kafka partition (optional, default is 0).                                 | `"partition"`           | Integer              | Yes         |
| `Kafka.TlsCertFile`       | The TLS certificate file path for Kafka (optional).                           | `"tls_cert_file"`       | String (File Path)   | Yes         |
| `Kafka.TlsKeyFile`        | The TLS key file path for Kafka (optional).                                   | `"tls_key_file"`        | String (File Path)   | Yes         |

### Webhook Object
| Field Name                | Description                                                                        | YAML Key                | Data Type            | Optional     |
|---------------------------|------------------------------------------------------------------------------------|-------------------------|----------------------|--------------|
| `Webhook`                 | Configuration for `Webhook` event source.                                          | `"webhook"`             | Object               | Yes          |
| `Webhook.Address`         | The Webhook address (required if `"webhook"` is in `EventSourceImpl`).             | `"address"`             | String               | Yes          |
| `Webhook.Path`            | The Webhook path (optional, must start with "/").                                  | `"path"`                | String               | Yes          |
| `Webhook.TlsCertFile`     | The TLS certificate file path for Webhook (required unless `TlsKeyFile` is empty). | `"tls_cert_file"`       | String (File Path)   | Yes          |
| `Webhook.TlsKeyFile`      | The TLS key file path for Webhook (required unless `TlsCertFile` is empty).        | `"tls_key_file"`        | String (File Path)   | Yes          |

### Ticker Object
| Field Name               | Description                                                | YAML Key       | Data Type            | Optional |
|--------------------------|------------------------------------------------------------|----------------|----------------------|----------|
| `Ticker`                 | Configuration for `Ticker` event source.                   | `"ticker"`     | Object               | Yes      |
| `Ticker.IntervalSeconds` | Seconds between invocations of the syncer. Must be >= 300. | `"interval_s"` | Integer (Seconds)    | No       |
