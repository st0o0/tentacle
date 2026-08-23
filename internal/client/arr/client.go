package arr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *Client) get(ctx context.Context, path string, query url.Values, target any) error {
	u := c.baseURL + "/api/v3" + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	var status SystemStatus
	if err := c.get(ctx, "/system/status", nil, &status); err != nil {
		return nil, fmt.Errorf("get system status: %w", err)
	}
	return &status, nil
}

func (c *Client) GetHealth(ctx context.Context) ([]HealthCheck, error) {
	var checks []HealthCheck
	if err := c.get(ctx, "/health", nil, &checks); err != nil {
		return nil, fmt.Errorf("get health: %w", err)
	}
	return checks, nil
}

func (c *Client) GetQueue(ctx context.Context) (*QueueResponse, error) {
	q := url.Values{}
	q.Set("page", "1")
	q.Set("pageSize", "1")

	var resp QueueResponse
	if err := c.get(ctx, "/queue", q, &resp); err != nil {
		return nil, fmt.Errorf("get queue: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	var folders []RootFolder
	if err := c.get(ctx, "/rootfolder", nil, &folders); err != nil {
		return nil, fmt.Errorf("get root folders: %w", err)
	}
	return folders, nil
}
