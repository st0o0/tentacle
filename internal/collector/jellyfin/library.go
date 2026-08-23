package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

var itemTypes = []string{
	"Movie", "Series", "Episode",
	"MusicAlbum", "MusicArtist", "Audio",
	"Book",
}

type LibraryCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	itemsTotal  *prometheus.Desc
	sizeBytes   *prometheus.Desc
	latestAdded *prometheus.Desc
}

func NewLibraryCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *LibraryCollector {
	return &LibraryCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		itemsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "items_total"),
			"Total number of items in a library by type.",
			[]string{"type", "library"}, nil,
		),
		sizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "size_bytes"),
			"Total size of a library in bytes.",
			[]string{"library"}, nil,
		),
		latestAdded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "latest_added_timestamp_seconds"),
			"Unix timestamp of the most recently added item in a library.",
			[]string{"library"}, nil,
		),
	}
}

func (c *LibraryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.itemsTotal
	ch <- c.sizeBytes
	ch <- c.latestAdded
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *LibraryCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	success := true

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	folders, err := c.client.GetVirtualFolders(ctx)
	if err != nil {
		c.logger.Error("library collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, time.Since(start).Seconds(), "library")
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "library")
		return
	}

	for _, folder := range folders {
		for _, itemType := range itemTypes {
			resp, err := c.client.GetItems(ctx, folder.ItemId, itemType, "", 0, 0)
			if err != nil {
				c.logger.Warn("failed to get item count", "library", folder.Name, "type", itemType, "err", err)
				success = false
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.itemsTotal, prometheus.GaugeValue, float64(resp.TotalRecordCount), itemType, folder.Name)
		}

		size, err := c.librarySize(ctx, folder.ItemId)
		if err != nil {
			c.logger.Warn("failed to get library size", "library", folder.Name, "err", err)
			success = false
		} else {
			ch <- prometheus.MustNewConstMetric(c.sizeBytes, prometheus.GaugeValue, float64(size), folder.Name)
		}

		latest, err := c.client.GetLatestItems(ctx, folder.ItemId, 1)
		if err != nil {
			c.logger.Warn("failed to get latest items", "library", folder.Name, "err", err)
			success = false
		} else if len(latest) > 0 {
			if t, ok := parseJellyfinTime(latest[0].DateCreated); ok {
				ch <- prometheus.MustNewConstMetric(c.latestAdded, prometheus.GaugeValue, float64(t.Unix()), folder.Name)
			}
		}
	}

	successVal := 1.0
	if !success {
		successVal = 0
	}
	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, time.Since(start).Seconds(), "library")
	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, successVal, "library")
}

func (c *LibraryCollector) librarySize(ctx context.Context, parentID string) (int64, error) {
	const batchSize = 500
	var total int64
	startIndex := 0

	for {
		resp, err := c.client.GetItems(ctx, parentID, "", "Size", batchSize, startIndex)
		if err != nil {
			return 0, err
		}
		for _, item := range resp.Items {
			total += item.Size
		}
		startIndex += len(resp.Items)
		if startIndex >= resp.TotalRecordCount || len(resp.Items) == 0 {
			break
		}
	}
	return total, nil
}
