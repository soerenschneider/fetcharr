package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/events"
	"github.com/soerenschneider/fetcharr/internal/events/kafka"
	"github.com/soerenschneider/fetcharr/internal/events/rabbitmq"
	"github.com/soerenschneider/fetcharr/internal/events/ticker"
	"github.com/soerenschneider/fetcharr/internal/events/webhook_server"
	"go.uber.org/multierr"
)

func buildEventSources(conf *config.Config) ([]events.EventSource, error) {
	var eventSources []events.EventSource

	var errs error
	for _, eventSourceImpl := range conf.EventSourceImpl {
		var err error
		var impl events.EventSource

		switch eventSourceImpl {
		case "kafka":
			impl, err = buildKafka(conf)
		case "webhook_server":
			impl, err = buildWebhook(conf)
		case "ticker":
			impl, err = buildTicker(conf)
		case "rabbitmq":
			impl, err = buildRabbitMq(conf)
		default:
			log.Warn().Msgf("Unknown event source impl: %s. This should not happen", eventSourceImpl)
		}

		if err != nil {
			errs = multierr.Append(errs, err)
		} else {
			eventSources = append(eventSources, impl)
		}
	}

	return eventSources, errs
}

func buildKafka(conf *config.Config) (*kafka.KafkaReader, error) {
	var opts []kafka.KafkaReaderOpts
	if conf.Kafka.Partition > 0 {
		opts = append(opts, kafka.WithPartition(conf.Kafka.Partition))
	}

	if len(conf.Kafka.TlsCertFile) > 0 && len(conf.Kafka.TlsKeyFile) > 0 {
		opts = append(opts, kafka.WithTlsCert(conf.Kafka.TlsCertFile))
		opts = append(opts, kafka.WithTlsKey(conf.Kafka.TlsKeyFile))
	}

	return kafka.NewReader(conf.Kafka.Brokers, conf.Kafka.Topic, conf.Kafka.GroupId, opts...)
}

func buildRabbitMq(conf *config.Config) (*rabbitmq.RabbitMqEventListener, error) {
	var opts []rabbitmq.RabbitMqOpts

	// add webhook_server path
	if len(conf.RabbitMq.ConsumerName) > 0 {
		opts = append(opts, rabbitmq.WithConsumerName(conf.RabbitMq.ConsumerName))
	}

	conn := rabbitmq.RabbitMqConnection{
		BrokerHost: conf.RabbitMq.Broker,
		Port:       conf.RabbitMq.Port,
		Username:   conf.RabbitMq.Username,
		Password:   conf.RabbitMq.Password,
		Vhost:      conf.RabbitMq.Vhost,
		CertFile:   conf.RabbitMq.TlsCertFile,
		KeyFile:    conf.RabbitMq.TlsKeyFile,
		UseSsl:     false,
	}

	return rabbitmq.New(conn, conf.RabbitMq.QueueName, opts...)
}

func buildWebhook(conf *config.Config) (*webhook_server.WebhookServer, error) {
	var opts []webhook_server.WebhookServerOpts

	// add webhook_server path
	if len(conf.Webhook.Path) > 0 {
		opts = append(opts, webhook_server.WithPath(conf.Webhook.Path))
	}

	// add tls keys
	if len(conf.Webhook.TlsCertFile) > 0 && len(conf.Webhook.TlsKeyFile) > 0 {
		opts = append(opts, webhook_server.WithTLS(conf.Webhook.TlsCertFile, conf.Webhook.TlsKeyFile))
	}

	return webhook_server.New(conf.Webhook.Address, opts...)
}

func buildTicker(conf *config.Config) (*ticker.Ticker, error) {
	interval := time.Second * time.Duration(conf.Ticker.IntervalSeconds)
	return ticker.NewTicker(interval)
}
