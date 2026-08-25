package prowlarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type systemCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	up          *prometheus.Desc
	info        *prometheus.Desc
	startTime   *prometheus.Desc
	healthTotal *prometheus.Desc
}

func newSystemCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *systemCollector {
	return &systemCollector{
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

func (c *systemCollector) Name() string { return "system" }

func (c *systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
	ch <- c.startTime
	ch <- c.healthTotal
}

func (c *systemCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status, err := c.client.GetSystemStatus(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return err
	}

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
		return err
	}

	counts := make(map[[2]string]int)
	for _, h := range health {
		counts[[2]string{h.Type, h.Source}]++
	}
	for key, count := range counts {
		ch <- prometheus.MustNewConstMetric(c.healthTotal, prometheus.GaugeValue, float64(count), key[0], key[1])
	}

	return nil
}
