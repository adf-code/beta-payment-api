package logger

import (
	"bytes"
	"github.com/adf-code/beta-payment-api/config"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"time"
)

type TelemetryClient struct {
	client   *http.Client
	apiKey   string
	endpoint string
}

func NewTelemetryClient(apiKey string, endpoint string) *TelemetryClient {
	return &TelemetryClient{
		client:   &http.Client{Timeout: 5 * time.Second},
		apiKey:   apiKey,
		endpoint: endpoint,
	}
}

func (w *TelemetryClient) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", w.endpoint, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.apiKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return len(p), nil
}

func InitLoggerWithTelemetry(cfg *config.AppConfig) zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	var writer io.Writer = consoleWriter

	if cfg.TelemetryEnabled == "true" {
		if cfg.TelemetryAPIKey == "" || cfg.TelemetryEndpoint == "" {
			panic("Telemetry config is not set")
		}

		telemetryWriter := NewTelemetryClient(cfg.TelemetryAPIKey, cfg.TelemetryEndpoint)
		writer = zerolog.MultiLevelWriter(consoleWriter, telemetryWriter)
	}

	return zerolog.New(writer).With().Timestamp().Logger()
}
