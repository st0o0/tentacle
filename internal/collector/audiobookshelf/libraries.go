package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type LibrariesCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	librariesTotal *prometheus.Desc
	itemsTotal     *prometheus.Desc
	sizeBytes      *prometheus.Desc
	duration       *prometheus.Desc
	audioTracks    *prometheus.Desc
	authorsTotal   *prometheus.Desc
	genresTotal    *prometheus.Desc
	missingTotal   *prometheus.Desc
	invalidTotal   *prometheus.Desc
	lastUpdate     *prometheus.Desc
}

func NewLibrariesCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *LibrariesCollector {
	labels := []string{"library", "media_type"}
	return &LibrariesCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		librariesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "libraries_total"),
			"Total number of libraries.",
			nil, nil,
		),
		itemsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "items_total"),
			"Total number of items in a library.",
			labels, nil,
		),
		sizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "size_bytes"),
			"Total size of a library in bytes.",
			labels, nil,
		),
		duration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "duration_seconds"),
			"Total duration of all items in a library in seconds.",
			labels, nil,
		),
		audioTracks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "audio_tracks_total"),
			"Total number of audio tracks in a library.",
			labels, nil,
		),
		authorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "authors_total"),
			"Total number of authors in a library.",
			labels, nil,
		),
		genresTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "genres_total"),
			"Total number of genres in a library.",
			labels, nil,
		),
		missingTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "missing_total"),
			"Total number of missing items in a library.",
			labels, nil,
		),
		invalidTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "invalid_total"),
			"Total number of invalid items in a library.",
			labels, nil,
		),
		lastUpdate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "last_update_timestamp_seconds"),
			"Unix timestamp of when a library was last updated.",
			labels, nil,
		),
	}
}

func (c *LibrariesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.librariesTotal
	ch <- c.itemsTotal
	ch <- c.sizeBytes
	ch <- c.duration
	ch <- c.audioTracks
	ch <- c.authorsTotal
	ch <- c.genresTotal
	ch <- c.missingTotal
	ch <- c.invalidTotal
	ch <- c.lastUpdate
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *LibrariesCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	libraries, err := c.client.GetLibraries(ctx)
	if err != nil {
		duration := time.Since(start).Seconds()
		c.logger.Error("libraries collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "libraries")
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "libraries")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.librariesTotal, prometheus.GaugeValue, float64(len(libraries)))

	for _, lib := range libraries {
		labels := []string{lib.Name, lib.MediaType}

		if lib.LastUpdate > 0 {
			ch <- prometheus.MustNewConstMetric(c.lastUpdate, prometheus.GaugeValue, float64(lib.LastUpdate)/1000.0, labels...)
		}

		stats, err := c.client.GetLibraryStats(ctx, lib.Id)
		if err != nil {
			c.logger.Error("library stats failed", "library", lib.Name, "err", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(c.itemsTotal, prometheus.GaugeValue, float64(stats.TotalItems), labels...)
		ch <- prometheus.MustNewConstMetric(c.sizeBytes, prometheus.GaugeValue, float64(stats.TotalSize), labels...)
		ch <- prometheus.MustNewConstMetric(c.duration, prometheus.GaugeValue, stats.TotalDuration, labels...)
		ch <- prometheus.MustNewConstMetric(c.audioTracks, prometheus.GaugeValue, float64(stats.NumAudioTracks), labels...)
		ch <- prometheus.MustNewConstMetric(c.authorsTotal, prometheus.GaugeValue, float64(stats.TotalAuthors), labels...)
		ch <- prometheus.MustNewConstMetric(c.genresTotal, prometheus.GaugeValue, float64(stats.TotalGenres), labels...)
		ch <- prometheus.MustNewConstMetric(c.missingTotal, prometheus.GaugeValue, float64(stats.NumMissing), labels...)
		ch <- prometheus.MustNewConstMetric(c.invalidTotal, prometheus.GaugeValue, float64(stats.NumInvalid), labels...)
	}

	duration := time.Since(start).Seconds()
	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "libraries")
	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "libraries")
}
