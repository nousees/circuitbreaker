package apiclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"circuit_breaker/config"

	"github.com/sony/gobreaker"
)

type ApiClient struct {
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
	baseURL        string
	maxRetries     int
	retryDelay     time.Duration
}

func NewApiClient(baseURL string, cfg *config.Config) *ApiClient {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        fmt.Sprintf("cb-%s", baseURL),
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.CBInterval,
		Timeout:     cfg.CBTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.FailureThreshold
		},
	})

	return &ApiClient{
		client:         &http.Client{Timeout: cfg.ClientTimeout},
		circuitBreaker: cb,
		baseURL:        baseURL,
		maxRetries:     cfg.MaxRetries,
		retryDelay:     cfg.RetryDelay,
	}
}
func (c *ApiClient) Call(ctx context.Context, endpoint string) (string, error) {
	url := c.baseURL + endpoint
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		body, err := c.circuitBreaker.Execute(func() (interface{}, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return nil, err
			}

			resp, err := c.client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 500 {
				return nil, fmt.Errorf("server error: %d", resp.StatusCode)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			return string(data), nil
		})

		if err == nil {
			return body.(string), nil
		}

		lastErr = err
		log.Printf("Attempt %d failed: %v", attempt+1, err)

		if attempt == c.maxRetries {
			break
		}

		time.Sleep(c.retryDelay * time.Duration(1<<attempt))
	}

	return "", fmt.Errorf("all retries failed: %w", lastErr)
}
