package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type backupsCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	backupsTotal *prometheus.Desc
	latestBackup *prometheus.Desc
}

func newBackupsCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *backupsCollector {
	return &backupsCollector{
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

func (c *backupsCollector) Name() string {
	return "backups"
}

func (c *backupsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.backupsTotal
	ch <- c.latestBackup
}

func (c *backupsCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetBackups(ctx)
	if err != nil {
		return err
	}

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

	return nil
}
