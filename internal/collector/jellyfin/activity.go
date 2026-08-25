package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type activityCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	totalEntries  *prometheus.Desc
	entriesByType *prometheus.Desc
	latestEntry   *prometheus.Desc
}

func newActivityCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *activityCollector {
	return &activityCollector{
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

func (c *activityCollector) Name() string { return "activity" }

func (c *activityCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalEntries
	ch <- c.entriesByType
	ch <- c.latestEntry
}

func (c *activityCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetActivityLog(ctx, 100)
	if err != nil {
		return err
	}

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
	return nil
}
