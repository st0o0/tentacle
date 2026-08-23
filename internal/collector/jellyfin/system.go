package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type SystemCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	info           *prometheus.Desc
	pendingRestart *prometheus.Desc
	up             *prometheus.Desc
}

func NewSystemCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *SystemCollector {
	return &SystemCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "info"),
			"Jellyfin system information.",
			[]string{"version", "os", "architecture"}, nil,
		),
		pendingRestart: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "pending_restart"),
			"Whether Jellyfin has a pending restart.",
			nil, nil,
		),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether Jellyfin is reachable.",
			nil, nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.info
	ch <- c.pendingRestart
	ch <- c.up
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	info, err := c.client.GetSystemInfo(ctx)
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
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, info.Version, info.OperatingSystem, info.SystemArchitecture)

	restart := 0.0
	if info.HasPendingRestart {
		restart = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.pendingRestart, prometheus.GaugeValue, restart)
}
