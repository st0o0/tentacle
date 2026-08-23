package arr

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

func (c *Client) GetSeries(ctx context.Context) ([]Series, error) {
	var series []Series
	if err := c.get(ctx, "/series", nil, &series); err != nil {
		return nil, fmt.Errorf("get series: %w", err)
	}
	return series, nil
}

func (c *Client) GetWantedMissing(ctx context.Context) (*WantedResponse, error) {
	q := url.Values{}
	q.Set("page", "1")
	q.Set("pageSize", "1")

	var resp WantedResponse
	if err := c.get(ctx, "/wanted/missing", q, &resp); err != nil {
		return nil, fmt.Errorf("get wanted missing: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetCalendar(ctx context.Context, start, end time.Time) ([]CalendarEntry, error) {
	q := url.Values{}
	q.Set("start", start.Format(time.RFC3339))
	q.Set("end", end.Format(time.RFC3339))

	var entries []CalendarEntry
	if err := c.get(ctx, "/calendar", q, &entries); err != nil {
		return nil, fmt.Errorf("get calendar: %w", err)
	}
	return entries, nil
}
