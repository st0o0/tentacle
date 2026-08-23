package prowlarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type IndexersCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	indexersTotal        *prometheus.Desc
	indexersEnabledTotal *prometheus.Desc
	queriesTotal         *prometheus.Desc
	grabsTotal           *prometheus.Desc
	failedQueriesTotal   *prometheus.Desc
	failedGrabsTotal     *prometheus.Desc
	avgResponseSeconds   *prometheus.Desc
}

func NewIndexersCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *IndexersCollector {
	return &IndexersCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		indexersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexers", "total"),
			"Total number of configured indexers.",
			nil, nil,
		),
		indexersEnabledTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexers", "enabled_total"),
			"Total number of enabled indexers.",
			nil, nil,
		),
		queriesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexer", "queries_total"),
			"Total number of queries per indexer.",
			[]string{"indexer"}, nil,
		),
		grabsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexer", "grabs_total"),
			"Total number of grabs per indexer.",
			[]string{"indexer"}, nil,
		),
		failedQueriesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexer", "failed_queries_total"),
			"Total number of failed queries per indexer.",
			[]string{"indexer"}, nil,
		),
		failedGrabsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexer", "failed_grabs_total"),
			"Total number of failed grabs per indexer.",
			[]string{"indexer"}, nil,
		),
		avgResponseSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexer", "avg_response_seconds"),
			"Average response time per indexer in seconds.",
			[]string{"indexer"}, nil,
		),
	}
}

func (c *IndexersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.indexersTotal
	ch <- c.indexersEnabledTotal
	ch <- c.queriesTotal
	ch <- c.grabsTotal
	ch <- c.failedQueriesTotal
	ch <- c.failedGrabsTotal
	ch <- c.avgResponseSeconds
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *IndexersCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	indexers, err := c.client.GetIndexers(ctx)
	if err != nil {
		duration := time.Since(start).Seconds()
		ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "indexers")
		c.logger.Error("indexers collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "indexers")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.indexersTotal, prometheus.GaugeValue, float64(len(indexers)))

	enabled := 0
	for _, idx := range indexers {
		if idx.Enable {
			enabled++
		}
	}
	ch <- prometheus.MustNewConstMetric(c.indexersEnabledTotal, prometheus.GaugeValue, float64(enabled))

	stats, err := c.client.GetIndexerStats(ctx)
	duration := time.Since(start).Seconds()
	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "indexers")

	if err != nil {
		c.logger.Error("indexer stats collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "indexers")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "indexers")

	for _, s := range stats.Indexers {
		ch <- prometheus.MustNewConstMetric(c.queriesTotal, prometheus.GaugeValue, float64(s.NumberOfQueries), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.grabsTotal, prometheus.GaugeValue, float64(s.NumberOfGrabs), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.failedQueriesTotal, prometheus.GaugeValue, float64(s.NumberOfFailedQueries), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.failedGrabsTotal, prometheus.GaugeValue, float64(s.NumberOfFailedGrabs), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.avgResponseSeconds, prometheus.GaugeValue, float64(s.AverageResponseTime)/1000.0, s.IndexerName)
	}
}
