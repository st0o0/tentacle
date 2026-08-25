package sonarr

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type extrasCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	backupTotal        *prometheus.Desc
	backupLatest       *prometheus.Desc
	updateAvailable    *prometheus.Desc
	blocklistTotal     *prometheus.Desc
	downloadClientInfo *prometheus.Desc
}

func newExtrasCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *extrasCollector {
	return &extrasCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		backupTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "backup", "total"),
			"Total number of backups.",
			nil, nil,
		),
		backupLatest: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "backup", "latest_timestamp_seconds"),
			"Unix timestamp of the most recent backup.",
			nil, nil,
		),
		updateAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "update", "available"),
			"Whether an update is available.",
			[]string{"version"}, nil,
		),
		blocklistTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "blocklist", "total"),
			"Total number of blocked releases.",
			nil, nil,
		),
		downloadClientInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "download_client", "info"),
			"Configured download client information.",
			[]string{"name", "protocol", "priority"}, nil,
		),
	}
}

func (c *extrasCollector) Name() string { return "extras" }

func (c *extrasCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.backupTotal
	ch <- c.backupLatest
	ch <- c.updateAvailable
	ch <- c.blocklistTotal
	ch <- c.downloadClientInfo
}

func (c *extrasCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var firstErr error

	backups, err := c.client.GetBackups(ctx)
	if err != nil {
		firstErr = err
	} else {
		ch <- prometheus.MustNewConstMetric(c.backupTotal, prometheus.GaugeValue, float64(len(backups)))
		var latestTime time.Time
		for _, b := range backups {
			if t, err := time.Parse(time.RFC3339, b.Time); err == nil && t.After(latestTime) {
				latestTime = t
			}
		}
		if !latestTime.IsZero() {
			ch <- prometheus.MustNewConstMetric(c.backupLatest, prometheus.GaugeValue, float64(latestTime.Unix()))
		}
	}

	updates, err := c.client.GetUpdates(ctx)
	if err != nil {
		if firstErr == nil {
			firstErr = err
		}
	} else if len(updates) > 0 {
		latest := updates[0]
		val := 0.0
		if !latest.Installed {
			val = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.updateAvailable, prometheus.GaugeValue, val, latest.Version)
	}

	blocklist, err := c.client.GetBlocklist(ctx)
	if err != nil {
		if firstErr == nil {
			firstErr = err
		}
	} else {
		ch <- prometheus.MustNewConstMetric(c.blocklistTotal, prometheus.GaugeValue, float64(blocklist.TotalRecords))
	}

	dlClients, err := c.client.GetDownloadClients(ctx)
	if err != nil {
		if firstErr == nil {
			firstErr = err
		}
	} else {
		for _, dc := range dlClients {
			if dc.Enable {
				ch <- prometheus.MustNewConstMetric(c.downloadClientInfo, prometheus.GaugeValue, 1, dc.Name, dc.Protocol, fmt.Sprintf("%d", dc.Priority))
			}
		}
	}

	return firstErr
}
