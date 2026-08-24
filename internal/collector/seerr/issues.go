package seerr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	client "github.com/st0o0/tentacle/internal/client/seerr"
)

type IssuesCollector struct {
	client  *client.Client
	timeout time.Duration
	logger  *slog.Logger

	issuesTotal    *prometheus.Desc
	issuesByStatus *prometheus.Desc
}

func NewIssuesCollector(c *client.Client, timeout time.Duration, logger *slog.Logger) *IssuesCollector {
	return &IssuesCollector{
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

func (c *IssuesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.issuesTotal
	ch <- c.issuesByStatus
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *IssuesCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	count, err := c.client.GetIssueCount(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "issues")

	if err != nil {
		c.logger.Error("issues collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "issues")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "issues")
	ch <- prometheus.MustNewConstMetric(c.issuesTotal, prometheus.GaugeValue, float64(count.Total))
	ch <- prometheus.MustNewConstMetric(c.issuesByStatus, prometheus.GaugeValue, float64(count.Open), "open")
	ch <- prometheus.MustNewConstMetric(c.issuesByStatus, prometheus.GaugeValue, float64(count.Resolved), "resolved")
}
