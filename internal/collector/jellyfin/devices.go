package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type devicesCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	total        *prometheus.Desc
	lastActivity *prometheus.Desc
}

func newDevicesCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *devicesCollector {
	return &devicesCollector{
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
			[]string{"device_name", "app_name", "app_version", "user"}, nil,
		),
	}
}

func (c *devicesCollector) Name() string { return "devices" }

func (c *devicesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.lastActivity
}

func (c *devicesCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetDevices(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(resp.TotalRecordCount))

	for _, d := range resp.Items {
		if t, ok := parseJellyfinTime(d.DateLastActivity); ok {
			ch <- prometheus.MustNewConstMetric(c.lastActivity, prometheus.GaugeValue, float64(t.Unix()), d.Name, d.AppName, d.AppVersion, d.LastUserName)
		}
	}
	return nil
}
