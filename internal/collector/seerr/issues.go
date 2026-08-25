package seerr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	client "github.com/st0o0/tentacle/internal/client/seerr"
)

type issuesCollector struct {
	client  *client.Client
	timeout time.Duration
	logger  *slog.Logger

	issuesTotal    *prometheus.Desc
	issuesByStatus *prometheus.Desc
}

func newIssuesCollector(c *client.Client, timeout time.Duration, logger *slog.Logger) *issuesCollector {
	return &issuesCollector{
		client:  c,
		timeout: timeout,
		logger:  logger,
		issuesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "issues", "total"),
			"Total number of reported issues.",
			nil, nil,
		),
		issuesByStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "issues", "by_status_total"),
			"Number of issues by status.",
			[]string{"status"}, nil,
		),
	}
}

func (c *issuesCollector) Name() string {
	return "issues"
}

func (c *issuesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.issuesTotal
	ch <- c.issuesByStatus
}

func (c *issuesCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	count, err := c.client.GetIssueCount(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.issuesTotal, prometheus.GaugeValue, float64(count.Total))
	ch <- prometheus.MustNewConstMetric(c.issuesByStatus, prometheus.GaugeValue, float64(count.Open), "open")
	ch <- prometheus.MustNewConstMetric(c.issuesByStatus, prometheus.GaugeValue, float64(count.Resolved), "resolved")

	return nil
}
