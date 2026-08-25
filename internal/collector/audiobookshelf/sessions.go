package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type sessionsCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	sessionsTotal *prometheus.Desc
}

func newSessionsCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *sessionsCollector {
	return &sessionsCollector{
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

func (c *sessionsCollector) Name() string {
	return "sessions"
}

func (c *sessionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.sessionsTotal
}

func (c *sessionsCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	sessions, err := c.client.GetSessions(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.sessionsTotal, prometheus.GaugeValue, float64(sessions.Total))

	return nil
}
