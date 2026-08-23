package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type SystemCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	up *prometheus.Desc
}

func NewSystemCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *SystemCollector {
	return &SystemCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether Audiobookshelf is reachable.",
			nil, nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	err := c.client.Ping(ctx)
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
}
