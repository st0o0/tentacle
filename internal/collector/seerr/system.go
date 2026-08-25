package seerr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/seerr"
)

type systemCollector struct {
	client  *seerr.Client
	timeout time.Duration
	logger  *slog.Logger

	info *prometheus.Desc
	up   *prometheus.Desc
}

func newSystemCollector(client *seerr.Client, timeout time.Duration, logger *slog.Logger) *systemCollector {
	return &systemCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "info"),
			"Seerr system information.",
			[]string{"version", "commit"}, nil,
		),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether Seerr is reachable.",
			nil, nil,
		),
	}
}

func (c *systemCollector) Name() string {
	return "system"
}

func (c *systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.info
	ch <- c.up
}

func (c *systemCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status, err := c.client.GetStatus(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, status.Version, status.CommitTag)

	return nil
}
