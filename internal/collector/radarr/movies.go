package radarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type moviesCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	total      *prometheus.Desc
	monitored  *prometheus.Desc
	byStatus   *prometheus.Desc
	downloaded *prometheus.Desc
	missing    *prometheus.Desc
	sizeBytes  *prometheus.Desc
}

func newMoviesCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *moviesCollector {
	return &moviesCollector{
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
		byStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "movies", "by_status_total"),
			"Number of movies by status.",
			[]string{"status"}, nil,
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

func (c *moviesCollector) Name() string { return "movies" }

func (c *moviesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.monitored
	ch <- c.byStatus
	ch <- c.downloaded
	ch <- c.missing
	ch <- c.sizeBytes
}

func (c *moviesCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	movies, err := c.client.GetMovies(ctx)
	if err != nil {
		return err
	}

	var monitoredCount, downloadedCount int
	var totalSize int64
	statusCounts := make(map[string]int)
	for _, m := range movies {
		if m.Monitored {
			monitoredCount++
		}
		if m.HasFile {
			downloadedCount++
		}
		totalSize += m.SizeOnDisk
		if m.Status != "" {
			statusCounts[m.Status]++
		}
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(len(movies)))
	ch <- prometheus.MustNewConstMetric(c.monitored, prometheus.GaugeValue, float64(monitoredCount))
	for status, count := range statusCounts {
		ch <- prometheus.MustNewConstMetric(c.byStatus, prometheus.GaugeValue, float64(count), status)
	}
	ch <- prometheus.MustNewConstMetric(c.downloaded, prometheus.GaugeValue, float64(downloadedCount))
	ch <- prometheus.MustNewConstMetric(c.sizeBytes, prometheus.GaugeValue, float64(totalSize))

	wanted, err := c.client.GetWantedMissingMovies(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.missing, prometheus.GaugeValue, float64(wanted.TotalRecords))

	return nil
}
