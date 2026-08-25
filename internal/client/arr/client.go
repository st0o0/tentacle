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
	apiVersion string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey, apiVersion string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		apiVersion: apiVersion,
		httpClient: httpClient,
	}
}

func (c *Client) get(ctx context.Context, path string, query url.Values, target any) error {
	u := c.baseURL + "/api/" + c.apiVersion + path
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
	q.Set("pageSize", "250")

	var resp QueueResponse
	if err := c.get(ctx, "/queue", q, &resp); err != nil {
		return nil, fmt.Errorf("get queue: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetBackups(ctx context.Context) ([]Backup, error) {
	var backups []Backup
	if err := c.get(ctx, "/system/backup", nil, &backups); err != nil {
		return nil, fmt.Errorf("get backups: %w", err)
	}
	return backups, nil
}

func (c *Client) GetUpdates(ctx context.Context) ([]Update, error) {
	var updates []Update
	if err := c.get(ctx, "/update", nil, &updates); err != nil {
		return nil, fmt.Errorf("get updates: %w", err)
	}
	return updates, nil
}

func (c *Client) GetBlocklist(ctx context.Context) (*BlocklistResponse, error) {
	q := url.Values{}
	q.Set("page", "1")
	q.Set("pageSize", "1")

	var resp BlocklistResponse
	if err := c.get(ctx, "/blocklist", q, &resp); err != nil {
		return nil, fmt.Errorf("get blocklist: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	var clients []DownloadClient
	if err := c.get(ctx, "/downloadclient", nil, &clients); err != nil {
		return nil, fmt.Errorf("get download clients: %w", err)
	}
	return clients, nil
}

func (c *Client) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	var folders []RootFolder
	if err := c.get(ctx, "/rootfolder", nil, &folders); err != nil {
		return nil, fmt.Errorf("get root folders: %w", err)
	}
	return folders, nil
}
