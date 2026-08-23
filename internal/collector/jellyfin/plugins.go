package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type PluginsCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	total           *prometheus.Desc
	info            *prometheus.Desc
	updateAvailable *prometheus.Desc
}

func NewPluginsCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *PluginsCollector {
	return &PluginsCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "plugins", "total"),
			"Total number of installed plugins.",
			nil, nil,
		),
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "plugin", "info"),
			"Installed plugin information.",
			[]string{"name", "version", "status"}, nil,
		),
		updateAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "plugin", "update_available"),
			"Whether a plugin has an update available.",
			[]string{"name"}, nil,
		),
	}
}

func (c *PluginsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.info
	ch <- c.updateAvailable
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *PluginsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	plugins, err := c.client.GetPlugins(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "plugins")

	if err != nil {
		c.logger.Error("plugins collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "plugins")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "plugins")
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(len(plugins)))

	for _, p := range plugins {
		ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, p.Name, p.Version, p.Status)

		update := 0.0
		if p.HasUpdate {
			update = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.updateAvailable, prometheus.GaugeValue, update, p.Name)
	}
}
