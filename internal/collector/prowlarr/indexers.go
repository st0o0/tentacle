package prowlarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type indexersCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	indexersTotal        *prometheus.Desc
	indexersEnabledTotal *prometheus.Desc
	indexersByProtocol   *prometheus.Desc
	queriesTotal         *prometheus.Desc
	grabsTotal           *prometheus.Desc
	failedQueriesTotal   *prometheus.Desc
	failedGrabsTotal     *prometheus.Desc
	avgResponseSeconds   *prometheus.Desc
}

func newIndexersCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *indexersCollector {
	return &indexersCollector{
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
		indexersByProtocol: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexers", "by_protocol_total"),
			"Number of indexers by protocol.",
			[]string{"protocol"}, nil,
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

func (c *indexersCollector) Name() string { return "indexers" }

func (c *indexersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.indexersTotal
	ch <- c.indexersEnabledTotal
	ch <- c.indexersByProtocol
	ch <- c.queriesTotal
	ch <- c.grabsTotal
	ch <- c.failedQueriesTotal
	ch <- c.failedGrabsTotal
	ch <- c.avgResponseSeconds
}

func (c *indexersCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	indexers, err := c.client.GetIndexers(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.indexersTotal, prometheus.GaugeValue, float64(len(indexers)))

	enabled := 0
	protocolCounts := make(map[string]int)
	for _, idx := range indexers {
		if idx.Enable {
			enabled++
		}
		if idx.Protocol != "" {
			protocolCounts[idx.Protocol]++
		}
	}
	ch <- prometheus.MustNewConstMetric(c.indexersEnabledTotal, prometheus.GaugeValue, float64(enabled))
	for proto, count := range protocolCounts {
		ch <- prometheus.MustNewConstMetric(c.indexersByProtocol, prometheus.GaugeValue, float64(count), proto)
	}

	stats, err := c.client.GetIndexerStats(ctx)
	if err != nil {
		return err
	}

	for _, s := range stats.Indexers {
		ch <- prometheus.MustNewConstMetric(c.queriesTotal, prometheus.GaugeValue, float64(s.NumberOfQueries), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.grabsTotal, prometheus.GaugeValue, float64(s.NumberOfGrabs), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.failedQueriesTotal, prometheus.GaugeValue, float64(s.NumberOfFailedQueries), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.failedGrabsTotal, prometheus.GaugeValue, float64(s.NumberOfFailedGrabs), s.IndexerName)
		ch <- prometheus.MustNewConstMetric(c.avgResponseSeconds, prometheus.GaugeValue, float64(s.AverageResponseTime)/1000.0, s.IndexerName)
	}

	return nil
}
