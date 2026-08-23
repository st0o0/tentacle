package seerr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/seerr"
)

type SystemCollector struct {
	client  *seerr.Client
	timeout time.Duration
	logger  *slog.Logger

	info *prometheus.Desc
	up   *prometheus.Desc
}

func NewSystemCollector(client *seerr.Client, timeout time.Duration, logger *slog.Logger) *SystemCollector {
	return &SystemCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "info"),
			"Seerr system information.",
			[]string{"version", "commit"}, nil,
		),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether Seerr is reachable.",
			nil, nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.info
	ch <- c.up
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status, err := c.client.GetStatus(ctx)
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
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, status.Version, status.CommitTag)
}
