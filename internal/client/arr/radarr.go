package arr

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

func (c *Client) GetMovies(ctx context.Context) ([]Movie, error) {
	var movies []Movie
	if err := c.get(ctx, "/movie", nil, &movies); err != nil {
		return nil, fmt.Errorf("get movies: %w", err)
	}
	return movies, nil
}

func (c *Client) GetWantedMissingMovies(ctx context.Context) (*WantedResponse, error) {
	q := url.Values{}
	q.Set("page", "1")
	q.Set("pageSize", "1")

	var resp WantedResponse
	if err := c.get(ctx, "/wanted/missing", q, &resp); err != nil {
		return nil, fmt.Errorf("get wanted missing movies: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetMovieCalendar(ctx context.Context, start, end time.Time) ([]CalendarEntry, error) {
	q := url.Values{}
	q.Set("start", start.Format(time.RFC3339))
	q.Set("end", end.Format(time.RFC3339))

	var entries []CalendarEntry
	if err := c.get(ctx, "/calendar", q, &entries); err != nil {
		return nil, fmt.Errorf("get movie calendar: %w", err)
	}
	return entries, nil
}
