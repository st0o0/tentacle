package seerr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/seerr"
)

type requestsCollector struct {
	client  *seerr.Client
	timeout time.Duration
	logger  *slog.Logger

	total    *prometheus.Desc
	byType   *prometheus.Desc
	byStatus *prometheus.Desc
}

func newRequestsCollector(client *seerr.Client, timeout time.Duration, logger *slog.Logger) *requestsCollector {
	return &requestsCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "requests", "total"),
			"Total number of Seerr requests.",
			nil, nil,
		),
		byType: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "requests_by_type", "total"),
			"Number of Seerr requests by media type.",
			[]string{"media_type"}, nil,
		),
		byStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "requests_by_status", "total"),
			"Number of Seerr requests by status.",
			[]string{"status"}, nil,
		),
	}
}

func (c *requestsCollector) Name() string {
	return "requests"
}

func (c *requestsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.byType
	ch <- c.byStatus
}

func (c *requestsCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	count, err := c.client.GetRequestCount(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(count.Total))
	ch <- prometheus.MustNewConstMetric(c.byType, prometheus.GaugeValue, float64(count.Movie), "movie")
	ch <- prometheus.MustNewConstMetric(c.byType, prometheus.GaugeValue, float64(count.TV), "tv")
	ch <- prometheus.MustNewConstMetric(c.byStatus, prometheus.GaugeValue, float64(count.Pending), "pending")
	ch <- prometheus.MustNewConstMetric(c.byStatus, prometheus.GaugeValue, float64(count.Approved), "approved")
	ch <- prometheus.MustNewConstMetric(c.byStatus, prometheus.GaugeValue, float64(count.Available), "available")
	ch <- prometheus.MustNewConstMetric(c.byStatus, prometheus.GaugeValue, float64(count.Declined), "declined")

	return nil
}
