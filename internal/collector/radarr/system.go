package radarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type systemCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	up           *prometheus.Desc
	info         *prometheus.Desc
	startTime    *prometheus.Desc
	healthIssues *prometheus.Desc
}

func newSystemCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *systemCollector {
	return &systemCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether Radarr is reachable.",
			nil, nil,
		),
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "info"),
			"Radarr system information.",
			[]string{"version", "branch", "runtime"}, nil,
		),
		startTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "start_time_seconds"),
			"Unix timestamp of when Radarr started.",
			nil, nil,
		),
		healthIssues: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "health", "issues_total"),
			"Number of health issues by type and source.",
			[]string{"type", "source"}, nil,
		),
	}
}

func (c *systemCollector) Name() string { return "system" }

func (c *systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
	ch <- c.startTime
	ch <- c.healthIssues
}

func (c *systemCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status, err := c.client.GetSystemStatus(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, status.Version, status.Branch, status.RuntimeName)

	if t, err := time.Parse(time.RFC3339, status.StartTime); err == nil {
		ch <- prometheus.MustNewConstMetric(c.startTime, prometheus.GaugeValue, float64(t.Unix()))
	}

	health, err := c.client.GetHealth(ctx)
	if err != nil {
		return err
	}

	for _, h := range health {
		ch <- prometheus.MustNewConstMetric(c.healthIssues, prometheus.GaugeValue, 1, h.Type, h.Source)
	}

	return nil
}
