package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type countsCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	itemsTotal *prometheus.Desc
}

func newCountsCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *countsCollector {
	return &countsCollector{
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

func (c *countsCollector) Name() string { return "counts" }

func (c *countsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.itemsTotal
}

func (c *countsCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	counts, err := c.client.GetItemCounts(ctx)
	if err != nil {
		return err
	}

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
	return nil
}
