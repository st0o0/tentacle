package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type systemCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	up   *prometheus.Desc
	info *prometheus.Desc
}

func newSystemCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *systemCollector {
	return &systemCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether Audiobookshelf is reachable.",
			nil, nil,
		),
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "info"),
			"Audiobookshelf system information.",
			[]string{"version"}, nil,
		),
	}
}

func (c *systemCollector) Name() string {
	return "system"
}

func (c *systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
}

func (c *systemCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	err := c.client.Ping(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

	backups, err := c.client.GetBackups(ctx)
	if err == nil && len(backups.Backups) > 0 {
		var latest audiobookshelf.Backup
		for _, b := range backups.Backups {
			if b.CreatedAt > latest.CreatedAt {
				latest = b
			}
		}
		if latest.ServerVersion != "" {
			ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, latest.ServerVersion)
		}
	}

	return nil
}
