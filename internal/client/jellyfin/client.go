package jellyfin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func (c *Client) get(ctx context.Context, path string, query url.Values, target any) error {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", c.token))

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

func (c *Client) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	var info SystemInfo
	if err := c.get(ctx, "/System/Info", nil, &info); err != nil {
		return nil, fmt.Errorf("get system info: %w", err)
	}
	return &info, nil
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	if err := c.get(ctx, "/Users", nil, &users); err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return users, nil
}

func (c *Client) GetSessions(ctx context.Context) ([]Session, error) {
	var sessions []Session
	if err := c.get(ctx, "/Sessions", nil, &sessions); err != nil {
		return nil, fmt.Errorf("get sessions: %w", err)
	}
	return sessions, nil
}

func (c *Client) GetVirtualFolders(ctx context.Context) ([]VirtualFolder, error) {
	var folders []VirtualFolder
	if err := c.get(ctx, "/Library/VirtualFolders", nil, &folders); err != nil {
		return nil, fmt.Errorf("get virtual folders: %w", err)
	}
	return folders, nil
}

func (c *Client) GetItemCounts(ctx context.Context) (*ItemCounts, error) {
	var counts ItemCounts
	if err := c.get(ctx, "/Items/Counts", nil, &counts); err != nil {
		return nil, fmt.Errorf("get item counts: %w", err)
	}
	return &counts, nil
}

func (c *Client) GetItems(ctx context.Context, parentID string, includeItemTypes string, fields string, limit, startIndex int) (*ItemsResponse, error) {
	q := url.Values{}
	if parentID != "" {
		q.Set("ParentId", parentID)
	}
	q.Set("Recursive", "true")
	if includeItemTypes != "" {
		q.Set("IncludeItemTypes", includeItemTypes)
	}
	if fields != "" {
		q.Set("Fields", fields)
	}
	q.Set("Limit", strconv.Itoa(limit))
	if startIndex > 0 {
		q.Set("StartIndex", strconv.Itoa(startIndex))
	}

	var resp ItemsResponse
	if err := c.get(ctx, "/Items", q, &resp); err != nil {
		return nil, fmt.Errorf("get items: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetLatestItems(ctx context.Context, parentID string, limit int) ([]Item, error) {
	q := url.Values{}
	if parentID != "" {
		q.Set("ParentId", parentID)
	}
	q.Set("Limit", strconv.Itoa(limit))

	var items []Item
	if err := c.get(ctx, "/Items/Latest", q, &items); err != nil {
		return nil, fmt.Errorf("get latest items: %w", err)
	}
	return items, nil
}

func (c *Client) GetScheduledTasks(ctx context.Context) ([]ScheduledTask, error) {
	var tasks []ScheduledTask
	if err := c.get(ctx, "/ScheduledTasks", nil, &tasks); err != nil {
		return nil, fmt.Errorf("get scheduled tasks: %w", err)
	}
	return tasks, nil
}

func (c *Client) GetActivityLog(ctx context.Context, limit int) (*ActivityLogResponse, error) {
	q := url.Values{}
	q.Set("Limit", strconv.Itoa(limit))

	var resp ActivityLogResponse
	if err := c.get(ctx, "/System/ActivityLog/Entries", q, &resp); err != nil {
		return nil, fmt.Errorf("get activity log: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetPlugins(ctx context.Context) ([]PluginInfo, error) {
	var plugins []PluginInfo
	if err := c.get(ctx, "/Plugins", nil, &plugins); err != nil {
		return nil, fmt.Errorf("get plugins: %w", err)
	}
	return plugins, nil
}

func (c *Client) GetDevices(ctx context.Context) (*DevicesResponse, error) {
	var resp DevicesResponse
	if err := c.get(ctx, "/Devices", nil, &resp); err != nil {
		return nil, fmt.Errorf("get devices: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetPlaybackActivity(ctx context.Context, days int) ([]PlaybackActivity, error) {
	q := url.Values{}
	q.Set("days", strconv.Itoa(days))

	var activity []PlaybackActivity
	if err := c.get(ctx, "/user_usage_stats/user_activity", q, &activity); err != nil {
		return nil, fmt.Errorf("get playback activity: %w", err)
	}
	return activity, nil
}
