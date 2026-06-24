package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL string, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) Get(ctx context.Context, path string, out any) error {
	body, err := c.request(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func (c *Client) Post(ctx context.Context, path string, payload any, out any) error {
	body, err := c.request(ctx, http.MethodPost, path, payload, map[string]string{"Prefer": "return=representation"})
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func (c *Client) Patch(ctx context.Context, path string, payload any, out any) error {
	body, err := c.request(ctx, http.MethodPatch, path, payload, map[string]string{"Prefer": "return=representation"})
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func (c *Client) request(ctx context.Context, method string, path string, payload any, headers map[string]string) ([]byte, error) {
	var reader io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase returned HTTP %d: %s", res.StatusCode, string(body))
	}
	return body, nil
}
