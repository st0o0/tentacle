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
