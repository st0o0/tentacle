package collector

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync"
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
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var opErr *net.OpError
	return errors.As(err, &opErr)
}

type SubCollector interface {
	Name() string
	Describe(ch chan<- *prometheus.Desc)
	Update(ch chan<- prometheus.Metric) error
}

// CachedCollector wraps a SubCollector and refreshes its metrics in the
// background on a fixed interval. During a Prometheus scrape it serves the
// most recently cached values instead of calling the upstream API.
type CachedCollector struct {
	inner    SubCollector
	interval time.Duration
	logger   *slog.Logger

	mu      sync.RWMutex
	cached  []prometheus.Metric
	lastErr error
	lastDur time.Duration
}

func NewCachedCollector(inner SubCollector, interval time.Duration, logger *slog.Logger) *CachedCollector {
	return &CachedCollector{
		inner:    inner,
		interval: interval,
		logger:   logger,
	}
}

func (c *CachedCollector) Name() string { return c.inner.Name() }

func (c *CachedCollector) Describe(ch chan<- *prometheus.Desc) {
	c.inner.Describe(ch)
}

func (c *CachedCollector) Update(ch chan<- prometheus.Metric) error {
	c.mu.RLock()
	metrics := c.cached
	err := c.lastErr
	c.mu.RUnlock()

	for _, m := range metrics {
		ch <- m
	}
	return err
}

// Start runs the background refresh loop. It performs an immediate first
// refresh, then refreshes on every tick until ctx is cancelled.
func (c *CachedCollector) Start(ctx context.Context) {
	c.refresh()

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.refresh()
		}
	}
}

func (c *CachedCollector) refresh() {
	ch := make(chan prometheus.Metric, 256)

	var wg sync.WaitGroup
	wg.Add(1)

	var metrics []prometheus.Metric
	go func() {
		defer wg.Done()
		for m := range ch {
			metrics = append(metrics, m)
		}
	}()

	start := time.Now()
	err := c.inner.Update(ch)
	close(ch)
	wg.Wait()
	dur := time.Since(start)

	c.mu.Lock()
	c.cached = metrics
	c.lastErr = err
	c.lastDur = dur
	c.mu.Unlock()

	if err != nil {
		c.logger.Warn("cached collector refresh failed", "collector", c.inner.Name(), "duration", dur.Seconds(), "err", err)
	} else {
		c.logger.Debug("cached collector refreshed", "collector", c.inner.Name(), "duration", dur.Seconds(), "metrics", len(metrics))
	}
}

type CacheIntervals struct {
	Warm time.Duration
	Cold time.Duration
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

// StartCaches launches background refresh goroutines for all
// CachedCollector instances among the registered sub-collectors.
func (sc *ServiceCollector) StartCaches(ctx context.Context) {
	for _, c := range sc.collectors {
		if cc, ok := c.(*CachedCollector); ok {
			go cc.Start(ctx)
		}
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
	serviceStart := time.Now()

	var wg sync.WaitGroup
	for _, c := range sc.collectors {
		wg.Add(1)
		go func(c SubCollector) {
			defer wg.Done()
			start := time.Now()
			err := c.Update(ch)
			duration := time.Since(start).Seconds()

			ch <- prometheus.MustNewConstMetric(sc.scrape.Duration, prometheus.GaugeValue, duration, c.Name())

			success := 1.0
			if err != nil {
				success = 0
				if isTransientError(err) {
					sc.logger.Warn("collector failed", "service", sc.namespace, "collector", c.Name(), "duration", duration, "err", err)
				} else {
					sc.logger.Error("collector failed", "service", sc.namespace, "collector", c.Name(), "duration", duration, "err", err)
				}
			}
			ch <- prometheus.MustNewConstMetric(sc.scrape.Success, prometheus.GaugeValue, success, c.Name())
		}(c)
	}
	wg.Wait()

	sc.logger.Debug("service collect complete", "service", sc.namespace, "duration", time.Since(serviceStart).Seconds())
}
