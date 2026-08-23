package audiobookshelf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		httpClient: httpClient,
	}
}

func (c *Client) get(ctx context.Context, path string, target any) error {
	u := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

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

func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/ping", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetStatus(ctx context.Context) (*ServerStatus, error) {
	var status ServerStatus
	if err := c.get(ctx, "/status", &status); err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}
	return &status, nil
}

func (c *Client) GetLibraries(ctx context.Context) ([]Library, error) {
	var resp struct {
		Libraries []Library `json:"libraries"`
	}
	if err := c.get(ctx, "/api/libraries", &resp); err != nil {
		return nil, fmt.Errorf("get libraries: %w", err)
	}
	return resp.Libraries, nil
}

func (c *Client) GetLibraryStats(ctx context.Context, libraryID string) (*LibraryStats, error) {
	var stats LibraryStats
	if err := c.get(ctx, "/api/libraries/"+libraryID+"/stats", &stats); err != nil {
		return nil, fmt.Errorf("get library stats: %w", err)
	}
	return &stats, nil
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	if err := c.get(ctx, "/api/users", &users); err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return users, nil
}

func (c *Client) GetOnlineUsers(ctx context.Context) ([]OnlineUser, error) {
	var resp struct {
		OpenSessions []ListeningSession `json:"openSessions"`
		Users        []OnlineUser       `json:"usersOnline"`
	}
	if err := c.get(ctx, "/api/users/online", &resp); err != nil {
		return nil, fmt.Errorf("get online users: %w", err)
	}
	return resp.Users, nil
}

func (c *Client) GetUserListeningStats(ctx context.Context, userID string) (*ListeningStats, error) {
	var stats ListeningStats
	if err := c.get(ctx, "/api/users/"+userID+"/listening-stats", &stats); err != nil {
		return nil, fmt.Errorf("get listening stats: %w", err)
	}
	return &stats, nil
}

func (c *Client) GetSessions(ctx context.Context) (*SessionsResponse, error) {
	var resp SessionsResponse
	if err := c.get(ctx, "/api/sessions", &resp); err != nil {
		return nil, fmt.Errorf("get sessions: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetBackups(ctx context.Context) (*BackupsResponse, error) {
	var resp BackupsResponse
	if err := c.get(ctx, "/api/backups", &resp); err != nil {
		return nil, fmt.Errorf("get backups: %w", err)
	}
	return &resp, nil
}
