package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type SessionsCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	sessionsTotal *prometheus.Desc
}

func NewSessionsCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *SessionsCollector {
	return &SessionsCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		sessionsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "sessions_total"),
			"Total number of listening sessions.",
			nil, nil,
		),
	}
}

func (c *SessionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.sessionsTotal
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SessionsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	sessions, err := c.client.GetSessions(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "sessions")

	if err != nil {
		c.logger.Error("sessions collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "sessions")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "sessions")
	ch <- prometheus.MustNewConstMetric(c.sessionsTotal, prometheus.GaugeValue, float64(sessions.Total))
}
