package sonarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type DiskCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	totalBytes *prometheus.Desc
	freeBytes  *prometheus.Desc
}

func NewDiskCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *DiskCollector {
	return &DiskCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		totalBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "disk", "total_bytes"),
			"Total disk space in bytes.",
			[]string{"path"}, nil,
		),
		freeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "disk", "free_bytes"),
			"Free disk space in bytes.",
			[]string{"path"}, nil,
		),
	}
}

func (c *DiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalBytes
	ch <- c.freeBytes
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *DiskCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	folders, err := c.client.GetRootFolders(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "disk")

	if err != nil {
		c.logger.Error("disk collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "disk")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "disk")

	for _, f := range folders {
		ch <- prometheus.MustNewConstMetric(c.totalBytes, prometheus.GaugeValue, float64(f.TotalSpace), f.Path)
		ch <- prometheus.MustNewConstMetric(c.freeBytes, prometheus.GaugeValue, float64(f.FreeSpace), f.Path)
	}
}
