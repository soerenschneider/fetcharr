# fetcharr
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/fetcharr)](https://goreportcard.com/report/github.com/soerenschneider/fetcharr)
![test-workflow](https://github.com/soerenschneider/fetcharr/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/fetcharr/actions/workflows/release-container.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/fetcharr/actions/workflows/golangci-lint.yaml/badge.svg)

## Features

🪆 A powerful wrapper around rsync (or any other command) to fetch data from remote systems<br/>
🔌 Multiple pluggable event notifiers that invoke sync process (Kafka, RabbitMQ, webhooks, time-based)<br/>
🪝 Support for defining multiple pre- and post-hooks<br/>
🔭 Observability through Prometheus metrics<br/>

## Roadmap

📣 Send notifications to user<br/>

## Why would I need this?

🔨 You want to fetch data from a remote system in a frequency that doesn't match a cron expression<br/>
📊 You want to get alerted on errors and look at dashboards rather than logs<br/>

## Installation

### Docker / Podman
````shell
$ git clone https://github.com/soerenschneider/fetcharr
$ cd fetcharr
$ docker run -v $(pwd)/contrib:/config ghcr.io/soerenschneider/fetcharr -config /config/fetcharr.yaml
````

### Binaries
Head over to the [prebuilt binaries](https://github.com/soerenschneider/fetcharr/releases) and download the correct binary for your system.
Use the example [systemd service file](contrib/fetcharr.service) to run it at boot.

### From Source
As a prerequesite, you need to have [Golang SDK](https://go.dev/dl/) installed. After that, you can install fetcharr from source by invoking:
```text
$ go install github.com/soerenschneider/fetcharr@latest
```

## Configuration
Head over to the [configuration section](docs/configuration.md) to see more details.

## Observability
Head over to the [metrics](docs/metrics.md) to see more details.

## Changelog
The changelog can be found [here](CHANGELOG.md)
