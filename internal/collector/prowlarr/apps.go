package prowlarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type AppsCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	appsTotal       *prometheus.Desc
	appInfo         *prometheus.Desc
	indexerDisabled *prometheus.Desc
}

func NewAppsCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *AppsCollector {
	return &AppsCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		appsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "apps", "total"),
			"Total number of connected applications.",
			nil, nil,
		),
		appInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "app", "info"),
			"Connected application information.",
			[]string{"name", "sync_level", "implementation"}, nil,
		),
		indexerDisabled: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "indexer", "disabled"),
			"Whether an indexer is temporarily disabled.",
			[]string{"indexer"}, nil,
		),
	}
}

func (c *AppsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.appsTotal
	ch <- c.appInfo
	ch <- c.indexerDisabled
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *AppsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	success := true

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	apps, err := c.client.GetApplications(ctx)
	if err != nil {
		c.logger.Error("applications fetch failed", "err", err)
		success = false
	} else {
		ch <- prometheus.MustNewConstMetric(c.appsTotal, prometheus.GaugeValue, float64(len(apps)))
		for _, app := range apps {
			ch <- prometheus.MustNewConstMetric(c.appInfo, prometheus.GaugeValue, 1, app.Name, app.SyncLevel, app.Implementation)
		}
	}

	indexers, err := c.client.GetIndexers(ctx)
	if err != nil {
		c.logger.Error("indexers fetch for status failed", "err", err)
		success = false
	} else {
		idToName := make(map[int]string)
		for _, idx := range indexers {
			idToName[idx.Id] = idx.Name
		}

		statuses, err := c.client.GetIndexerStatuses(ctx)
		if err != nil {
			c.logger.Error("indexer statuses fetch failed", "err", err)
			success = false
		} else {
			now := time.Now()
			for _, s := range statuses {
				name := idToName[s.IndexerId]
				if name == "" {
					continue
				}
				val := 0.0
				if t, err := time.Parse(time.RFC3339, s.DisabledTill); err == nil && t.After(now) {
					val = 1.0
				}
				ch <- prometheus.MustNewConstMetric(c.indexerDisabled, prometheus.GaugeValue, val, name)
			}
		}
	}

	successVal := 1.0
	if !success {
		successVal = 0
	}
	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, time.Since(start).Seconds(), "apps")
	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, successVal, "apps")
}
