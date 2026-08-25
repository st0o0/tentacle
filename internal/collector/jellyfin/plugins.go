package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type pluginsCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	total           *prometheus.Desc
	info            *prometheus.Desc
	updateAvailable *prometheus.Desc
}

func newPluginsCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *pluginsCollector {
	return &pluginsCollector{
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

func (c *pluginsCollector) Name() string { return "plugins" }

func (c *pluginsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.info
	ch <- c.updateAvailable
}

func (c *pluginsCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	plugins, err := c.client.GetPlugins(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(len(plugins)))

	for _, p := range plugins {
		ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, p.Name, p.Version, p.Status)

		update := 0.0
		if p.HasUpdate {
			update = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.updateAvailable, prometheus.GaugeValue, update, p.Name)
	}
	return nil
}
