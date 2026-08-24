package arr

import (
	"context"
	"fmt"
)

func (c *Client) GetIndexers(ctx context.Context) ([]Indexer, error) {
	var indexers []Indexer
	if err := c.get(ctx, "/indexer", nil, &indexers); err != nil {
		return nil, fmt.Errorf("get indexers: %w", err)
	}
	return indexers, nil
}

func (c *Client) GetIndexerStats(ctx context.Context) (*IndexerStatsResponse, error) {
	var resp IndexerStatsResponse
	if err := c.get(ctx, "/indexerstats", nil, &resp); err != nil {
		return nil, fmt.Errorf("get indexer stats: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetApplications(ctx context.Context) ([]Application, error) {
	var apps []Application
	if err := c.get(ctx, "/applications", nil, &apps); err != nil {
		return nil, fmt.Errorf("get applications: %w", err)
	}
	return apps, nil
}

func (c *Client) GetIndexerStatuses(ctx context.Context) ([]IndexerStatus, error) {
	var statuses []IndexerStatus
	if err := c.get(ctx, "/indexerstatus", nil, &statuses); err != nil {
		return nil, fmt.Errorf("get indexer statuses: %w", err)
	}
	return statuses, nil
}
