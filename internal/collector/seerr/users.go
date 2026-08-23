package seerr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/seerr"
)

type UsersCollector struct {
	client  *seerr.Client
	timeout time.Duration
	logger  *slog.Logger

	total *prometheus.Desc
}

func NewUsersCollector(client *seerr.Client, timeout time.Duration, logger *slog.Logger) *UsersCollector {
	return &UsersCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "users", "total"),
			"Total number of Seerr users.",
			nil, nil,
		),
	}
}

func (c *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *UsersCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetUsers(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "users")

	if err != nil {
		c.logger.Error("users collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "users")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "users")
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(resp.PageInfo.Results))
}
