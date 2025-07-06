package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/soerenschneider/fetcharr/internal/events"
	"go.uber.org/multierr"
)

const defaultTimeout = 10 * time.Second

type KafkaReader struct {
	brokers   []string
	topic     string
	groupId   string
	partition int

	tlsKey  string
	tlsCert string
}

type KafkaReaderOpts func(*KafkaReader) error

func NewReader(brokers []string, topic string, groupId string, opts ...KafkaReaderOpts) (*KafkaReader, error) {
	if len(brokers) == 0 {
		return nil, errors.New("empty list of kafka brokers supplied")
	}

	if len(topic) == 0 {
		return nil, errors.New("empty topic supplied")
	}

	if len(groupId) == 0 {
		return nil, errors.New("empty groupId supplied")
	}

	kafka := &KafkaReader{
		topic:   topic,
		brokers: brokers,
		groupId: groupId,
	}

	var errs error
	for _, opt := range opts {
		err := opt(kafka)
		errs = multierr.Append(errs, err)
	}

	return kafka, errs
}

func (k *KafkaReader) Listen(ctx context.Context, eventChan chan events.EventSyncRequest, wg *sync.WaitGroup) error {
	if ctx == nil {
		return errors.New("empty context supplied")
	}

	if eventChan == nil {
		return errors.New("closed channel supplied")
	}

	wg.Add(1)

	dialer := &kafka.Dialer{
		Timeout:   defaultTimeout,
		DualStack: true,
		TLS: &tls.Config{
			GetClientCertificate: k.loadTlsClientCerts,
			MinVersion:           tls.VersionTLS12,
		},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.brokers,
		Topic:     k.topic,
		Partition: k.partition,
		GroupID:   k.groupId,
		MaxBytes:  10e6,
		Dialer:    dialer,
	})

	continueConsuming := true
	for continueConsuming {
		select {
		case <-ctx.Done():
			log.Info().Msg("Kafka received signal")
			continueConsuming = false
		default:
			_, err := r.ReadMessage(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Error().Err(err).Msg("Error while reading kafka message")
			} else {
				req := events.EventSyncRequest{
					Source:   "kafka",
					Metadata: "", // TODO
					Response: make(chan error),
				}

				eventChan <- req

				select {
				case err := <-req.Response:
					if err != nil {
						log.Error().Str("component", "kafka").Err(err).Msg("received error response from fetcharr")
					}
				case <-time.After(3 * time.Second):
					log.Warn().Str("component", "kafka").Msgf("timeout waiting for goroutine")
				}
			}
		}
	}

	err := r.Close()
	wg.Done()
	return err
}

func (k *KafkaReader) loadTlsClientCerts(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	if len(k.tlsCert) == 0 || len(k.tlsKey) == 0 {
		return nil, errors.New("no client certificates defined")
	}

	certificate, err := tls.LoadX509KeyPair(k.tlsCert, k.tlsKey)
	if err != nil {
		log.Error().Err(err).Msg("user-defined client certificates could not be loaded")
	}
	return &certificate, err
}
