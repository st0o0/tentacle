package radarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type queueCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	total   *prometheus.Desc
	byState *prometheus.Desc
}

func newQueueCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *queueCollector {
	return &queueCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "queue", "total"),
			"Total number of items in the queue.",
			nil, nil,
		),
		byState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "queue", "by_state_total"),
			"Number of queue items by tracked download state.",
			[]string{"state"}, nil,
		),
	}
}

func (c *queueCollector) Name() string { return "queue" }

func (c *queueCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.byState
}

func (c *queueCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	queue, err := c.client.GetQueue(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(queue.TotalRecords))

	stateCounts := make(map[string]int)
	for _, r := range queue.Records {
		if r.TrackedDownloadState != "" {
			stateCounts[r.TrackedDownloadState]++
		}
	}
	for state, count := range stateCounts {
		ch <- prometheus.MustNewConstMetric(c.byState, prometheus.GaugeValue, float64(count), state)
	}

	return nil
}
