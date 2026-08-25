package seerr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/seerr"
)

type usersCollector struct {
	client  *seerr.Client
	timeout time.Duration
	logger  *slog.Logger

	total *prometheus.Desc
}

func newUsersCollector(client *seerr.Client, timeout time.Duration, logger *slog.Logger) *usersCollector {
	return &usersCollector{
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

func (c *usersCollector) Name() string {
	return "users"
}

func (c *usersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
}

func (c *usersCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.GetUsers(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(resp.PageInfo.Results))

	return nil
}
