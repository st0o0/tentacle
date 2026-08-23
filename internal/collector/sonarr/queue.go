package sonarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type QueueCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	queueTotal *prometheus.Desc
}

func NewQueueCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *QueueCollector {
	return &QueueCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		queueTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "queue", "total"),
			"Total number of items in the download queue.",
			nil, nil,
		),
	}
}

func (c *QueueCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.queueTotal
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *QueueCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	queue, err := c.client.GetQueue(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "queue")

	if err != nil {
		c.logger.Error("queue collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "queue")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "queue")
	ch <- prometheus.MustNewConstMetric(c.queueTotal, prometheus.GaugeValue, float64(queue.TotalRecords))
}
