package seerr

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
	u := c.baseURL + "/api/v1" + path
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

func (c *Client) GetStatus(ctx context.Context) (*Status, error) {
	var status Status
	if err := c.get(ctx, "/status", nil, &status); err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}
	return &status, nil
}

func (c *Client) GetRequestCount(ctx context.Context) (*RequestCount, error) {
	var count RequestCount
	if err := c.get(ctx, "/request/count", nil, &count); err != nil {
		return nil, fmt.Errorf("get request count: %w", err)
	}
	return &count, nil
}

func (c *Client) GetIssueCount(ctx context.Context) (*IssueCount, error) {
	var count IssueCount
	if err := c.get(ctx, "/issue/count", nil, &count); err != nil {
		return nil, fmt.Errorf("get issue count: %w", err)
	}
	return &count, nil
}

func (c *Client) GetUsers(ctx context.Context) (*UsersResponse, error) {
	q := url.Values{}
	q.Set("take", "100")
	q.Set("skip", "0")

	var resp UsersResponse
	if err := c.get(ctx, "/user", q, &resp); err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return &resp, nil
}
