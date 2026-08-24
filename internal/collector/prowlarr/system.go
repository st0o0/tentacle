package prowlarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type SystemCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	up          *prometheus.Desc
	info        *prometheus.Desc
	startTime   *prometheus.Desc
	healthTotal *prometheus.Desc
}

func NewSystemCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *SystemCollector {
	return &SystemCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Whether Prowlarr is reachable.",
			nil, nil,
		),
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "info"),
			"Prowlarr system information.",
			[]string{"version", "branch", "runtime"}, nil,
		),
		startTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "start_time_seconds"),
			"Unix timestamp of when Prowlarr started.",
			nil, nil,
		),
		healthTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "health", "issues_total"),
			"Number of health issues by type and source.",
			[]string{"type", "source"}, nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
	ch <- c.startTime
	ch <- c.healthTotal
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status, err := c.client.GetSystemStatus(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "system")

	if err != nil {
		c.logger.Error("system collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "system")
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "system")
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

	runtime := status.RuntimeName
	if status.RuntimeVersion != "" {
		runtime += " " + status.RuntimeVersion
	}
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, status.Version, status.Branch, runtime)

	if t, err := time.Parse(time.RFC3339, status.StartTime); err == nil {
		ch <- prometheus.MustNewConstMetric(c.startTime, prometheus.GaugeValue, float64(t.Unix()))
	}

	health, err := c.client.GetHealth(ctx)
	if err != nil {
		c.logger.Error("health collector failed", "err", err)
		return
	}

	counts := make(map[[2]string]int)
	for _, h := range health {
		counts[[2]string{h.Type, h.Source}]++
	}
	for key, count := range counts {
		ch <- prometheus.MustNewConstMetric(c.healthTotal, prometheus.GaugeValue, float64(count), key[0], key[1])
	}
}
