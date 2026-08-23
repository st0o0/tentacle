package radarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type MoviesCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	total       *prometheus.Desc
	monitored   *prometheus.Desc
	downloaded  *prometheus.Desc
	missing     *prometheus.Desc
	sizeBytes   *prometheus.Desc
}

func NewMoviesCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *MoviesCollector {
	return &MoviesCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "movies", "total"),
			"Total number of movies.",
			nil, nil,
		),
		monitored: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "movies", "monitored_total"),
			"Total number of monitored movies.",
			nil, nil,
		),
		downloaded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "movies", "downloaded_total"),
			"Total number of downloaded movies.",
			nil, nil,
		),
		missing: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "movies", "missing_total"),
			"Total number of missing movies.",
			nil, nil,
		),
		sizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "movies", "size_bytes"),
			"Total size of all movies in bytes.",
			nil, nil,
		),
	}
}

func (c *MoviesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.monitored
	ch <- c.downloaded
	ch <- c.missing
	ch <- c.sizeBytes
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *MoviesCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	movies, err := c.client.GetMovies(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "movies")

	if err != nil {
		c.logger.Error("movies collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "movies")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "movies")

	var monitoredCount, downloadedCount int
	var totalSize int64
	for _, m := range movies {
		if m.Monitored {
			monitoredCount++
		}
		if m.HasFile {
			downloadedCount++
		}
		totalSize += m.SizeOnDisk
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(len(movies)))
	ch <- prometheus.MustNewConstMetric(c.monitored, prometheus.GaugeValue, float64(monitoredCount))
	ch <- prometheus.MustNewConstMetric(c.downloaded, prometheus.GaugeValue, float64(downloadedCount))
	ch <- prometheus.MustNewConstMetric(c.sizeBytes, prometheus.GaugeValue, float64(totalSize))

	wanted, err := c.client.GetWantedMissingMovies(ctx)
	if err != nil {
		c.logger.Error("wanted missing movies failed", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.missing, prometheus.GaugeValue, float64(wanted.TotalRecords))
}
