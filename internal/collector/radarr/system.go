package radarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type SystemCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	up           *prometheus.Desc
	info         *prometheus.Desc
	healthIssues *prometheus.Desc
}

func NewSystemCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *SystemCollector {
	return &SystemCollector{
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
		healthIssues: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "health", "issues_total"),
			"Number of health issues by type and source.",
			[]string{"type", "source"}, nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
	ch <- c.healthIssues
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status, err := c.client.GetSystemStatus(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "system")

	if err != nil {
		c.logger.Error("system collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "system")
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "system")
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, status.Version, status.Branch, status.RuntimeName)

	health, err := c.client.GetHealth(ctx)
	if err != nil {
		c.logger.Error("health check failed", "err", err)
		return
	}

	for _, h := range health {
		ch <- prometheus.MustNewConstMetric(c.healthIssues, prometheus.GaugeValue, 1, h.Type, h.Source)
	}
}
