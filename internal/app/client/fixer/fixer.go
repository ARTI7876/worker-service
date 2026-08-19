package fixer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ARTI7876/worker-service/internal/app/config/section"
	"github.com/ARTI7876/worker-service/internal/app/entity"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const requestTimeout = 10 * time.Second

// Client — HTTP-клиент Fixer API. Транспорт обёрнут otelhttp, поэтому
// исходящие запросы трассируются (no-op при выключенном OpenTelemetry).
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type fixerResponse struct {
	Success bool               `json:"success"`
	Base    string             `json:"base"`
	Date    string             `json:"date"`
	Rates   map[string]float64 `json:"rates"`
	Error   *fixerError        `json:"error,omitempty"`
}

type fixerError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

func NewClient(cfg section.ClientFixer) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   requestTimeout,
		},
		apiKey:  cfg.ApiKey,
		baseURL: cfg.BaseURL,
	}
}

// GetRates возвращает курсы валют относительно базовой валюты base.
func (c *Client) GetRates(ctx context.Context, base string) (map[string]float64, error) {
	args := url.Values{}
	args.Set("access_key", c.apiKey)
	args.Set("base", base)

	requestURL := fmt.Sprintf("%s/latest?%s", c.baseURL, args.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrFixerUnavailable, err.Error())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrFixerUnavailable, err.Error())
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", entity.ErrFixerUnavailable, resp.StatusCode)
	}

	var body fixerResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrFixerInvalidResponse, err.Error())
	}

	if !body.Success {
		return nil, mapFixerError(body.Error)
	}

	log.Info().
		Ctx(ctx).
		Str("base", base).
		Int("count", len(body.Rates)).
		Msg("Получены курсы валют от Fixer")

	return body.Rates, nil
}

// mapFixerError преобразует ошибку Fixer API в типизированную (errors.Is-совместимую).
func mapFixerError(fixerErr *fixerError) error {
	if fixerErr == nil {
		return entity.ErrFixerInvalidResponse
	}

	switch fixerErr.Code {
	case 101:
		return entity.ErrFixerInvalidApiKey
	case 104, 105:
		return entity.ErrFixerRateLimitExceeded
	default:
		return fmt.Errorf(
			"%w: [%d] %s - %s",
			entity.ErrFixerInvalidResponse,
			fixerErr.Code,
			fixerErr.Type,
			fixerErr.Info,
		)
	}
}
