package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type systemCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	info           *prometheus.Desc
	pendingRestart *prometheus.Desc
	up             *prometheus.Desc
}

func newSystemCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *systemCollector {
	return &systemCollector{
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

func (c *systemCollector) Name() string { return "system" }

func (c *systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.info
	ch <- c.pendingRestart
	ch <- c.up
}

func (c *systemCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	info, err := c.client.GetSystemInfo(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, info.Version, info.OperatingSystem, info.SystemArchitecture)

	restart := 0.0
	if info.HasPendingRestart {
		restart = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.pendingRestart, prometheus.GaugeValue, restart)
	return nil
}
