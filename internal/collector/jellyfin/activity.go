package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type ActivityCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	totalEntries   *prometheus.Desc
	entriesByType  *prometheus.Desc
	latestEntry    *prometheus.Desc
}

func NewActivityCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *ActivityCollector {
	return &ActivityCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		totalEntries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "activity_log", "entries_total"),
			"Total number of activity log entries by severity.",
			[]string{"severity"}, nil,
		),
		entriesByType: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "activity_log", "entries_by_type_total"),
			"Total number of activity log entries by type.",
			[]string{"type"}, nil,
		),
		latestEntry: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "activity_log", "latest_timestamp_seconds"),
			"Unix timestamp of the most recent activity log entry.",
			nil, nil,
		),
	}
}

func (c *ActivityCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalEntries
	ch <- c.entriesByType
	ch <- c.latestEntry
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *ActivityCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetActivityLog(ctx, 100)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "activity")

	if err != nil {
		c.logger.Error("activity collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "activity")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "activity")

	severityCounts := map[string]float64{
		"Information": 0,
		"Warning":     0,
		"Error":       0,
	}
	typeCounts := make(map[string]float64)

	var latestTime time.Time

	for _, entry := range resp.Items {
		if entry.Severity != "" {
			severityCounts[entry.Severity]++
		}
		if entry.Type != "" {
			typeCounts[entry.Type]++
		}
		if t, ok := parseJellyfinTime(entry.Date); ok && t.After(latestTime) {
			latestTime = t
		}
	}

	for severity, count := range severityCounts {
		ch <- prometheus.MustNewConstMetric(c.totalEntries, prometheus.GaugeValue, count, severity)
	}
	for entryType, count := range typeCounts {
		ch <- prometheus.MustNewConstMetric(c.entriesByType, prometheus.GaugeValue, count, entryType)
	}

	if !latestTime.IsZero() {
		ch <- prometheus.MustNewConstMetric(c.latestEntry, prometheus.GaugeValue, float64(latestTime.Unix()))
	}
}
