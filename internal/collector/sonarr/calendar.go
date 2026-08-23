package sonarr

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

	upcomingTotal *prometheus.Desc
}

func NewCalendarCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *CalendarCollector {
	return &CalendarCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		upcomingTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "calendar", "upcoming_total"),
			"Total number of upcoming episodes in the next 7 days.",
			nil, nil,
		),
	}
}

func (c *CalendarCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upcomingTotal
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *CalendarCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	now := time.Now()
	entries, err := c.client.GetCalendar(ctx, now, now.Add(7*24*time.Hour))
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "calendar")

	if err != nil {
		c.logger.Error("calendar collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "calendar")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "calendar")
	ch <- prometheus.MustNewConstMetric(c.upcomingTotal, prometheus.GaugeValue, float64(len(entries)))
}
