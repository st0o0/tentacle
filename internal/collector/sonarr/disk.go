package sonarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type diskCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	totalBytes *prometheus.Desc
	freeBytes  *prometheus.Desc
}

func newDiskCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *diskCollector {
	return &diskCollector{
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

func (c *diskCollector) Name() string { return "disk" }

func (c *diskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalBytes
	ch <- c.freeBytes
}

func (c *diskCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	folders, err := c.client.GetRootFolders(ctx)
	if err != nil {
		return err
	}

	for _, f := range folders {
		ch <- prometheus.MustNewConstMetric(c.totalBytes, prometheus.GaugeValue, float64(f.TotalSpace), f.Path)
		ch <- prometheus.MustNewConstMetric(c.freeBytes, prometheus.GaugeValue, float64(f.FreeSpace), f.Path)
	}

	return nil
}
