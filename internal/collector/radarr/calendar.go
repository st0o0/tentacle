package radarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type CalendarCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	upcoming *prometheus.Desc
}

func NewCalendarCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *CalendarCollector {
	return &CalendarCollector{
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

func (c *CalendarCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upcoming
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *CalendarCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	now := time.Now()
	entries, err := c.client.GetMovieCalendar(ctx, now, now.AddDate(0, 0, 30))
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "calendar")

	if err != nil {
		c.logger.Error("calendar collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "calendar")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "calendar")
	ch <- prometheus.MustNewConstMetric(c.upcoming, prometheus.GaugeValue, float64(len(entries)))
}
