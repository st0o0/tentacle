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

type libraryCollector struct {
	client    *jellyfin.Client
	timeout   time.Duration
	batchSize int
	logger    *slog.Logger

	itemsTotal  *prometheus.Desc
	sizeBytes   *prometheus.Desc
	latestAdded *prometheus.Desc
}

func newLibraryCollector(client *jellyfin.Client, timeout time.Duration, batchSize int, logger *slog.Logger) *libraryCollector {
	return &libraryCollector{
		client:    client,
		timeout:   timeout,
		batchSize: batchSize,
		logger:    logger,
		itemsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "items_total"),
			"Total number of items in a library by type.",
			[]string{"type", "library", "collection_type"}, nil,
		),
		sizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "size_bytes"),
			"Total size of a library in bytes.",
			[]string{"library", "collection_type"}, nil,
		),
		latestAdded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "latest_added_timestamp_seconds"),
			"Unix timestamp of the most recently added item in a library.",
			[]string{"library", "collection_type"}, nil,
		),
	}
}

func (c *libraryCollector) Name() string { return "library" }

func (c *libraryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.itemsTotal
	ch <- c.sizeBytes
	ch <- c.latestAdded
}

func (c *libraryCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	folders, err := c.client.GetVirtualFolders(ctx)
	if err != nil {
		return err
	}

	for _, folder := range folders {
		ct := folder.CollectionType
		for _, itemType := range itemTypes {
			resp, err := c.client.GetItems(ctx, folder.ItemId, itemType, "", 0, 0)
			if err != nil {
				c.logger.Warn("failed to get item count", "library", folder.Name, "type", itemType, "err", err)
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.itemsTotal, prometheus.GaugeValue, float64(resp.TotalRecordCount), itemType, folder.Name, ct)
		}

		size, err := c.librarySize(ctx, folder.ItemId)
		if err != nil {
			c.logger.Warn("failed to get library size", "library", folder.Name, "err", err)
		} else {
			ch <- prometheus.MustNewConstMetric(c.sizeBytes, prometheus.GaugeValue, float64(size), folder.Name, ct)
		}

		latest, err := c.client.GetLatestItems(ctx, folder.ItemId, 1)
		if err != nil {
			c.logger.Warn("failed to get latest items", "library", folder.Name, "err", err)
		} else if len(latest) > 0 {
			if t, ok := parseJellyfinTime(latest[0].DateCreated); ok {
				ch <- prometheus.MustNewConstMetric(c.latestAdded, prometheus.GaugeValue, float64(t.Unix()), folder.Name, ct)
			}
		}
	}

	return nil
}

func (c *libraryCollector) librarySize(ctx context.Context, parentID string) (int64, error) {
	var total int64
	startIndex := 0

	for {
		resp, err := c.client.GetItems(ctx, parentID, "Movie,Episode,Audio,MusicVideo,Book", "MediaSources", c.batchSize, startIndex)
		if err != nil {
			return 0, err
		}
		for _, item := range resp.Items {
			for _, ms := range item.MediaSources {
				total += ms.Size
			}
		}
		startIndex += len(resp.Items)
		if startIndex >= resp.TotalRecordCount || len(resp.Items) == 0 {
			break
		}
	}
	return total, nil
}
