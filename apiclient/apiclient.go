package apiclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type ApiClient struct {
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
	baseURL        string
	maxRetries     int
	retryDelay     time.Duration
}

func NewApiClient(baseURL string, maxRetries int, failureThreshold uint32) *ApiClient {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        fmt.Sprintf("cb-%s", baseURL),
		MaxRequests: 2,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= failureThreshold
		},
	})

	return &ApiClient{
		client:         &http.Client{Timeout: 5 * time.Second},
		circuitBreaker: cb,
		baseURL:        baseURL,
		maxRetries:     maxRetries,
		retryDelay:     500 * time.Millisecond,
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
