package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/events"
	"github.com/soerenschneider/fetcharr/internal/metrics"
	"go.uber.org/multierr"
)

type RabbitMqConnection struct {
	BrokerHost string
	Port       int
	Username   string
	Password   string
	Vhost      string
	UseSsl     bool

	CertFile string
	KeyFile  string
}

type RabbitMqEventListener struct {
	connection   RabbitMqConnection
	queueName    string
	consumerName string
}

type RabbitMqOpts func(listener *RabbitMqEventListener) error

func New(conn RabbitMqConnection, queueName string, opts ...RabbitMqOpts) (*RabbitMqEventListener, error) {
	ret := &RabbitMqEventListener{
		connection:   conn,
		queueName:    queueName,
		consumerName: "fetcharr",
	}

	var errs error
	for _, opt := range opts {
		if err := opt(ret); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return ret, errs
}

func (e *RabbitMqEventListener) buildConnectionString() string {
	protocol := "amqp"
	if e.connection.UseSsl {
		protocol = "amqps"
	}

	if len(e.connection.CertFile) > 0 && len(e.connection.KeyFile) > 0 {
		return fmt.Sprintf("amqps://%s:%d/%s", e.connection.BrokerHost, e.connection.Port, e.connection.Vhost)
	}

	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", protocol, e.connection.Username, e.connection.Password, e.connection.BrokerHost, e.connection.Port, e.connection.Vhost)
}

func (e *RabbitMqEventListener) Listen(ctx context.Context, eventChan chan events.EventSyncRequest, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	operation := func() (any, error) {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			err := e.listen(ctx, eventChan)
			if err != nil {
				log.Error().Str("component", "rabbitmq").Err(err).Msg("error while listening on rabbitmq events")
				amqpErr := &amqp.Error{}
				if errors.As(err, amqpErr) {
					metrics.RabbitMqErrors.WithLabelValues(strconv.Itoa(amqpErr.Code)).Inc()
				}
			}
			return nil, err
		}
	}

	cont := true
	for cont {
		select {
		case <-ctx.Done():
			log.Debug().Str("component", "rabbitmq").Msg("Packing up")
			cont = false
		default:
			_, err := backoff.Retry[any](ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
			if err != nil {
				return fmt.Errorf("too many errors trying to listen on rabbitmq: %w", err)
			}
		}
	}

	return nil
}

// nolint: cyclop
func (e *RabbitMqEventListener) listen(ctx context.Context, eventChan chan events.EventSyncRequest) error {
	conn, err := amqp.Dial(e.buildConnectionString())
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()
	conNotify := conn.NotifyClose(make(chan *amqp.Error, 1))

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	chNotify := ch.NotifyClose(make(chan *amqp.Error, 1))

	msgs, err := ch.Consume(
		e.queueName,
		e.consumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return backoff.Permanent(err)
	}

	log.Info().Str("component", "rabbitmq").Msg("Listening for messages...")
	for {
		select {
		case err := <-conNotify:
			log.Warn().Err(err).Str("component", "rabbitmq").Msg("connection closed")
			metrics.RabbitMqDisconnects.WithLabelValues("connection").Inc()
			return err
		case err := <-chNotify:
			log.Warn().Err(err).Str("component", "rabbitmq").Msg("channel closed")
			metrics.RabbitMqDisconnects.WithLabelValues("channel").Inc()
			return err
		case <-msgs:
			log.Debug().Str("component", "rabbitmq").Msg("received message")
			req := events.EventSyncRequest{
				Source:   "rabbitmq",
				Metadata: "", // TODO
				Response: make(chan error),
			}

			eventChan <- req

			select {
			case err := <-req.Response:
				if err != nil {
					log.Error().Str("component", "rabbitmq").Err(err).Msg("received error response from fetcharr")
				}
			case <-time.After(3 * time.Second):
				log.Warn().Str("component", "rabbitmq").Msgf("timeout waiting for goroutine")
			}
		case <-ctx.Done():
			log.Debug().Str("component", "rabbitmq").Msg("context done")
			return nil
		}
	}
}
