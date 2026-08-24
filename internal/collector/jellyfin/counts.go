package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type CountsCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	itemsTotal *prometheus.Desc
}

func NewCountsCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *CountsCollector {
	return &CountsCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		itemsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "items", "total"),
			"Global item count by media type.",
			[]string{"type"}, nil,
		),
	}
}

func (c *CountsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.itemsTotal
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *CountsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	counts, err := c.client.GetItemCounts(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "counts")

	if err != nil {
		c.logger.Error("counts collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "counts")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "counts")

	emit := func(mediaType string, count int) {
		ch <- prometheus.MustNewConstMetric(c.itemsTotal, prometheus.GaugeValue, float64(count), mediaType)
	}

	emit("Movie", counts.MovieCount)
	emit("Series", counts.SeriesCount)
	emit("Episode", counts.EpisodeCount)
	emit("MusicArtist", counts.ArtistCount)
	emit("MusicAlbum", counts.AlbumCount)
	emit("Audio", counts.SongCount)
	emit("Book", counts.BookCount)
	emit("MusicVideo", counts.MusicVideoCount)
	emit("Trailer", counts.TrailerCount)
	emit("BoxSet", counts.BoxSetCount)
	emit("Program", counts.ProgramCount)
	emit("Item", counts.ItemCount)
}
