package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/fetcharr/internal/config"
	"github.com/soerenschneider/fetcharr/internal/syncer"
	"github.com/soerenschneider/fetcharr/pkg"
)

type WebhookClient struct {
	client *http.Client

	conf config.WebhookConfig
}

type WebhookClientOpts func(client *WebhookClient) error

func NewWebhookClient(client *http.Client, conf config.WebhookConfig) (*WebhookClient, error) {
	w := &WebhookClient{
		client: client,
		conf:   conf,
	}

	return w, nil
}

func (w *WebhookClient) Run(ctx context.Context, stats *syncer.Stats) error {
	data, err := w.getData(stats)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, w.conf.Verb, w.conf.Endpoint, data)
	if err != nil {
		return err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}

	if err := resp.Body.Close(); err != nil {
		log.Warn().Err(err).Msg("could not close body, this is weird and should not happen")
	}

	if resp.StatusCode >= 400 {
		// Do not print endpoint as the path may be considered sensitive information
		return fmt.Errorf("invoking webhook %q produced status code %d", w.conf.Name, resp.StatusCode)
	}

	return err
}

func (w *WebhookClient) getData(stats *syncer.Stats) (io.Reader, error) {
	if w.conf.Data != nil {
		jsonBody, err := json.Marshal(w.conf.Data)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(jsonBody), nil
	}

	if len(w.conf.EncodedDataFile) > 0 && stats != nil {
		return getTemplatedPayload(w.conf.EncodedDataFile, *stats)
	}

	return nil, nil
}

func getTemplatedPayload(templateFile string, stats syncer.Stats) (io.Reader, error) {
	data, err := os.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	templated, err := pkg.Format(string(data), stats)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(templated), nil
}

func (w *WebhookClient) GetType() string {
	return w.conf.GetType()
}

func (w *WebhookClient) GetName() string {
	return w.conf.GetName()
}

func (w *WebhookClient) GetStage() config.Stage {
	return w.conf.GetStage()
}
