package radarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type calendarCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	upcoming *prometheus.Desc
}

func newCalendarCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *calendarCollector {
	return &calendarCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		upcoming: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "calendar", "upcoming_total"),
			"Total number of upcoming movies in the next 30 days.",
			nil, nil,
		),
	}
}

func (c *calendarCollector) Name() string { return "calendar" }

func (c *calendarCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upcoming
}

func (c *calendarCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	now := time.Now()
	entries, err := c.client.GetMovieCalendar(ctx, now, now.AddDate(0, 0, 30))
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.upcoming, prometheus.GaugeValue, float64(len(entries)))

	return nil
}
