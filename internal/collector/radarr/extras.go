package radarr

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type ExtrasCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	backupTotal        *prometheus.Desc
	backupLatest       *prometheus.Desc
	updateAvailable    *prometheus.Desc
	blocklistTotal     *prometheus.Desc
	downloadClientInfo *prometheus.Desc
}

func NewExtrasCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *ExtrasCollector {
	return &ExtrasCollector{
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

func (c *ExtrasCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.backupTotal
	ch <- c.backupLatest
	ch <- c.updateAvailable
	ch <- c.blocklistTotal
	ch <- c.downloadClientInfo
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *ExtrasCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	success := true

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	backups, err := c.client.GetBackups(ctx)
	if err != nil {
		c.logger.Error("backups fetch failed", "err", err)
		success = false
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
		c.logger.Error("updates fetch failed", "err", err)
		success = false
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
		c.logger.Error("blocklist fetch failed", "err", err)
		success = false
	} else {
		ch <- prometheus.MustNewConstMetric(c.blocklistTotal, prometheus.GaugeValue, float64(blocklist.TotalRecords))
	}

	dlClients, err := c.client.GetDownloadClients(ctx)
	if err != nil {
		c.logger.Error("download clients fetch failed", "err", err)
		success = false
	} else {
		for _, dc := range dlClients {
			if dc.Enable {
				ch <- prometheus.MustNewConstMetric(c.downloadClientInfo, prometheus.GaugeValue, 1, dc.Name, dc.Protocol, fmt.Sprintf("%d", dc.Priority))
			}
		}
	}

	successVal := 1.0
	if !success {
		successVal = 0
	}
	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, time.Since(start).Seconds(), "extras")
	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, successVal, "extras")
}
