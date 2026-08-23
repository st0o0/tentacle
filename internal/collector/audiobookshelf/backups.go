package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type BackupsCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	backupsTotal   *prometheus.Desc
	latestBackup   *prometheus.Desc
}

func NewBackupsCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *BackupsCollector {
	return &BackupsCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		backupsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "backups_total"),
			"Total number of backups.",
			nil, nil,
		),
		latestBackup: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "backup", "latest_timestamp_seconds"),
			"Timestamp of the latest backup in seconds.",
			nil, nil,
		),
	}
}

func (c *BackupsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.backupsTotal
	ch <- c.latestBackup
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *BackupsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetBackups(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "backups")

	if err != nil {
		c.logger.Error("backups collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "backups")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "backups")
	ch <- prometheus.MustNewConstMetric(c.backupsTotal, prometheus.GaugeValue, float64(len(resp.Backups)))

	if len(resp.Backups) > 0 {
		var maxCreatedAt int64
		for _, b := range resp.Backups {
			if b.CreatedAt > maxCreatedAt {
				maxCreatedAt = b.CreatedAt
			}
		}
		ch <- prometheus.MustNewConstMetric(c.latestBackup, prometheus.GaugeValue, float64(maxCreatedAt)/1000.0)
	}
}
