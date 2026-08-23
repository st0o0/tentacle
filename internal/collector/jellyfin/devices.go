package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type DevicesCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	total        *prometheus.Desc
	lastActivity *prometheus.Desc
}

func NewDevicesCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *DevicesCollector {
	return &DevicesCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "devices", "total"),
			"Total number of registered devices.",
			nil, nil,
		),
		lastActivity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "device", "last_activity_timestamp_seconds"),
			"Unix timestamp of a device's last activity.",
			[]string{"device_name", "app_name", "user"}, nil,
		),
	}
}

func (c *DevicesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.lastActivity
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *DevicesCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetDevices(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "devices")

	if err != nil {
		c.logger.Error("devices collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "devices")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "devices")
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(resp.TotalRecordCount))

	for _, d := range resp.Items {
		if t, ok := parseJellyfinTime(d.DateLastActivity); ok {
			ch <- prometheus.MustNewConstMetric(c.lastActivity, prometheus.GaugeValue, float64(t.Unix()), d.Name, d.AppName, d.LastUserName)
		}
	}
}
