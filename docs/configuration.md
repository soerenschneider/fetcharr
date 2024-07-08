# Configuration

## Config Reference
| Field Name        | Description                                                                                 | YAML Key         | Data Type            | Optional |
|-------------------|---------------------------------------------------------------------------------------------|------------------|----------------------|----------|
| `SyncerImpl`      | The syncer implementation (currently only "rsync" is supported)                             | `"syncer_impl"`  | String               | No       |
| `Hooks`           | Defines hooks that are supposed to be run at certain stages                                 | `"hooks"`        | Array of Hooks       | Yes      |
| `EventSourceImpl` | An array of event source implementations. Valid choices are "kafka", "webhook" and "ticker" | `"events_impl"`  | Array of strings     | Yes      |
| `MetricsAddr`     | The metrics server address (optional, must be in the format "host:port").                   | `"metrics_addr"` | String (TCP Address) | Yes      |
| `Rsync`           | Configuration for the Rsync syncer                                                          | `"rsync"`        | Object               | No       |
| `Kafka`           | Kafka event source configuration                                                            | `"kafka"`        | Kafka (see below)    | Yes      |
| `RabbitMQ`        | RabbitMQ event source configuration                                                         | `"rabbitmq"`     | RabbitMq (see below) | Yes      |
| `Webhook`         | Webhook event source configuration                                                          | `"webhook"`      | Webhook (see below)  | Yes      |
| `Ticker`          | Ticker event source configuration                                                           | `"ticker"`       | Ticker (see below)   | Yes      |

### Rsync Object
| Field Name          | Description                                                          | YAML Key                | Data Type            | Optional    |
|---------------------|----------------------------------------------------------------------|-------------------------|----------------------|-------------|
| `Host`              | The Rsync host.                                                      | `"host"`                | String               | Yes         |
| `LocalDir`          | The local directory (must end with "/").                             | `"local_dir"`           | String               | No          |
| `RemoteDir`         | The remote directory (must end with "/").                            | `"remote_dir"`          | String               | No          |
| `BandwidthLimit`    | The bandwidth limit for                                              | `"bwlimit"`             | String               | Yes         |
| `ExcludePattern`    | The pattern for excluding files during synchronization.              | `"exclude"`             | String               | Yes         |
| `RemoveSourceFiles` | Flag to remove source files after successful sync.                   | `"remove_source_files"` | Boolean              | Yes         |
| `RemoteShell`       | The remote shell command (optional). Can be used to set ssh options. | `"remote_shell"`        | String               | Yes         |

## Event Sources

### Kafka Object
| Field Name         | Description                                                                  | YAML Key                | Data Type            | Optional |
|--------------------|------------------------------------------------------------------------------|-------------------------|----------------------|----------|
| `Brokers`          | The list of Kafka brokers (required if `"kafka"` is in `EventSourceImpl`).   | `"brokers"`             | Array of strings     | No       |
| `Topic`            | The Kafka topic (required if `"kafka"` is in `EventSourceImpl`).             | `"topic"`               | String               | No       |
| `GroupId`          | The Kafka consumer group ID (required if `"kafka"` is in `EventSourceImpl`). | `"group_id"`            | String               | No       |
| `Partition`        | The Kafka partition (optional, default is 0).                                | `"partition"`           | Integer              | Yes      |
| `TlsCertFile`      | The TLS certificate file path for Kafka (optional).                          | `"tls_cert_file"`       | String (File Path)   | Yes      |
| `TlsKeyFile`       | The TLS key file path for Kafka (optional).                                  | `"tls_key_file"`        | String (File Path)   | Yes      |

### RabbitMq Object
| Field Name     | Description                                                                   | YAML Key          | Data Type          | Optional |
|----------------|-------------------------------------------------------------------------------|-------------------|--------------------|----------|
| `Broker`       | The list of Kafka brokers (required if `"kafka"` is in `EventSourceImpl`).    | `"broker"`        | string             | No       |
| `Port`         | The Kafka topic (required if `"kafka"` is in `EventSourceImpl`).              | `"port"`          | Integer            | No       |
| `QueueName`    | The Kafka consumer group ID (required if `"kafka"` is in `EventSourceImpl`).  | `"queue"`         | String             | No       |
| `Vhost`        | The Kafka partition (optional, default is 0).                                 | `"vhost"`         | String             | No       |
| `ConsumerName` | The Kafka partition (optional, default is 0).                                 | `"consumer"`      | String             | Yes      |
| `Username`     | The Kafka partition (optional, default is 0).                                 | `"username"`      | String             | No       |
| `Password`     | The Kafka partition (optional, default is 0).                                 | `"password"`      | String             | No       |
| `TlsCertFile`  | The TLS certificate file path for Kafka (optional).                           | `"tls_cert_file"` | String (File Path) | Yes      |
| `TlsKeyFile`   | The TLS key file path for Kafka (optional).                                   | `"tls_key_file"`  | String (File Path) | Yes      |

### Webhook Event Source Object
| Field Name      | Description                                                                        | YAML Key                | Data Type            | Optional     |
|-----------------|------------------------------------------------------------------------------------|-------------------------|----------------------|--------------|
| `Address`       | The Webhook address (required if `"webhook"` is in `EventSourceImpl`).             | `"address"`             | String               | Yes          |
| `Path`          | The Webhook path (optional, must start with "/").                                  | `"path"`                | String               | Yes          |
| `TlsCertFile`   | The TLS certificate file path for Webhook (required unless `TlsKeyFile` is empty). | `"tls_cert_file"`       | String (File Path)   | Yes          |
| `TlsKeyFile`    | The TLS key file path for Webhook (required unless `TlsCertFile` is empty).        | `"tls_key_file"`        | String (File Path)   | Yes          |

### Ticker Object
| Field Name         | Description                                                | YAML Key       | Data Type            | Optional |
|--------------------|------------------------------------------------------------|----------------|----------------------|----------|
| `IntervalSeconds`  | Seconds between invocations of the syncer. Must be >= 300. | `"interval_s"` | Integer (Seconds)    | No       |

## Hooks

### Command Hook Object
| Field Name        | Description                                                                                                       | YAML Key              | Data Type          | Optional |
|-------------------|-------------------------------------------------------------------------------------------------------------------|-----------------------|--------------------|----------|
| `Name`            | The name of the hook. Used in logs and in metrics.                                                                | `"name"`              | String             | No       |
| `Stage`           | Hooks can run at different stages. This variable defines at which stage it should run.                            | `"stage"`             | Stage Enum         | No       |
| `Cmds`            | The command with its arguments that is supposed to be run                                                         | `"cmds"`              | Array of string    | No       |
| `ExitOnError`     | Whether a failure in this hook should be ignored or stop the application. Can only be used if the stage is `PRE`. | `"exit_on_error"`     | Bool               | Yes      |

###  Webhook Object
| Field Name        | Description                                                                                                       | YAML Key              | Data Type          | Optional |
|-------------------|-------------------------------------------------------------------------------------------------------------------|-----------------------|--------------------|----------|
| `Name`            | The name of the hook. Used in logs and in metrics.                                                                | `"name"`              | String             | No       |
| `Stage`           | Hooks can run at different stages. This variable defines at which stage it should run.                            | `"stage"`             | Stage Enum         | No       |
| `Endpoint`        | The endpoint to connect to                                                                                        | `"endpoint"`          | String             | No       |
| `Verb`            | The HTTP verb to use. Must be one of GET, POST, PATCH, HEAD.                                                      | `"verb"`              | String             | No       |
| `Data`            | The data to send to the endpoint                                                                                  | `"data"`              | String             | Yes      |
| `EncodedDataFile` | Path of a local file's content to send to the endpoint                                                            | `"encoded_data_file"` | String (File Path) | Yes      |
| `ExitOnError`     | Whether a failure in this hook should be ignored or stop the application. Can only be used if the stage is `PRE`. | `"exit_on_error"`     | Bool               | Yes      |

## Stages

### Stage Enum
| Field Name              | When is it executed?                                                            |
|-------------------------|---------------------------------------------------------------------------------|
| `PRE`                   | Before a sync is started                                                        |
| `POST_SUCCESS`          | After a successful sync happened (possibly no actual data has been transferred) |
| `POST_SUCCESS_TRANSFER` | After a successful sync happened where actual data has been transferred         |
| `POST_FAILURE`          | After a sync that was not successful                                            |
