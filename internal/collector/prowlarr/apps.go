package prowlarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type appsCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	appsTotal       *prometheus.Desc
	appInfo         *prometheus.Desc
	indexerDisabled *prometheus.Desc
}

func newAppsCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *appsCollector {
	return &appsCollector{
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

func (c *appsCollector) Name() string { return "apps" }

func (c *appsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.appsTotal
	ch <- c.appInfo
	ch <- c.indexerDisabled
}

func (c *appsCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var firstErr error

	apps, err := c.client.GetApplications(ctx)
	if err != nil {
		firstErr = err
	} else {
		ch <- prometheus.MustNewConstMetric(c.appsTotal, prometheus.GaugeValue, float64(len(apps)))
		for _, app := range apps {
			ch <- prometheus.MustNewConstMetric(c.appInfo, prometheus.GaugeValue, 1, app.Name, app.SyncLevel, app.Implementation)
		}
	}

	indexers, err := c.client.GetIndexers(ctx)
	if err != nil {
		if firstErr == nil {
			firstErr = err
		}
	} else {
		idToName := make(map[int]string)
		for _, idx := range indexers {
			idToName[idx.Id] = idx.Name
		}

		statuses, err := c.client.GetIndexerStatuses(ctx)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
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

	return firstErr
}
