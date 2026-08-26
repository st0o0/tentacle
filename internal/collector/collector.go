package collector

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client"
)

func isTransientError(err error) bool {
	if client.IsStatusCode(err, http.StatusBadGateway) ||
		client.IsStatusCode(err, http.StatusServiceUnavailable) ||
		client.IsStatusCode(err, http.StatusGatewayTimeout) {
		return true
	}
	return errors.Is(err, context.DeadlineExceeded)
}

type SubCollector interface {
	Name() string
	Describe(ch chan<- *prometheus.Desc)
	Update(ch chan<- prometheus.Metric) error
}

type ScrapeDescs struct {
	Duration *prometheus.Desc
	Success  *prometheus.Desc
}

func NewScrapeDescs(namespace string) ScrapeDescs {
	return ScrapeDescs{
		Duration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "duration_seconds"),
			"Duration of a collector scrape.",
			[]string{"collector"}, nil,
		),
		Success: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "success"),
			"Whether a collector scrape was successful.",
			[]string{"collector"}, nil,
		),
	}
}

type ServiceCollector struct {
	namespace  string
	scrape     ScrapeDescs
	logger     *slog.Logger
	collectors []SubCollector
}

func NewServiceCollector(namespace string, logger *slog.Logger, subs ...SubCollector) *ServiceCollector {
	return &ServiceCollector{
		namespace:  namespace,
		scrape:     NewScrapeDescs(namespace),
		logger:     logger,
		collectors: subs,
	}
}

func (sc *ServiceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sc.scrape.Duration
	ch <- sc.scrape.Success
	for _, c := range sc.collectors {
		c.Describe(ch)
	}
}

func (sc *ServiceCollector) Collect(ch chan<- prometheus.Metric) {
	for _, c := range sc.collectors {
		start := time.Now()
		err := c.Update(ch)
		duration := time.Since(start).Seconds()

		ch <- prometheus.MustNewConstMetric(sc.scrape.Duration, prometheus.GaugeValue, duration, c.Name())

		success := 1.0
		if err != nil {
			success = 0
			if isTransientError(err) {
				sc.logger.Warn("collector failed", "service", sc.namespace, "collector", c.Name(), "err", err)
			} else {
				sc.logger.Error("collector failed", "service", sc.namespace, "collector", c.Name(), "err", err)
			}
		}
		ch <- prometheus.MustNewConstMetric(sc.scrape.Success, prometheus.GaugeValue, success, c.Name())
	}
}
