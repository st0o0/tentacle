package sonarr

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
)

type SeriesCollector struct {
	client  *arr.Client
	timeout time.Duration
	logger  *slog.Logger

	seriesTotal     *prometheus.Desc
	monitored       *prometheus.Desc
	seriesByStatus  *prometheus.Desc
	seasonsTotal    *prometheus.Desc
	episodesTotal   *prometheus.Desc
	downloaded      *prometheus.Desc
	missing         *prometheus.Desc
	seriesSizeBytes *prometheus.Desc
}

func NewSeriesCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *SeriesCollector {
	return &SeriesCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		seriesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "series", "total"),
			"Total number of series.",
			nil, nil,
		),
		monitored: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "series", "monitored_total"),
			"Total number of monitored series.",
			nil, nil,
		),
		seriesByStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "series", "by_status_total"),
			"Number of series by status.",
			[]string{"status"}, nil,
		),
		seasonsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "seasons", "total"),
			"Total number of seasons across all series.",
			nil, nil,
		),
		episodesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "episodes", "total"),
			"Total number of episodes across all series.",
			nil, nil,
		),
		downloaded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "episodes", "downloaded_total"),
			"Total number of downloaded episodes.",
			nil, nil,
		),
		missing: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "episodes", "missing_total"),
			"Total number of missing episodes.",
			nil, nil,
		),
		seriesSizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "series", "size_bytes"),
			"Total size of all series on disk in bytes.",
			nil, nil,
		),
	}
}

func (c *SeriesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.seriesTotal
	ch <- c.monitored
	ch <- c.seriesByStatus
	ch <- c.seasonsTotal
	ch <- c.episodesTotal
	ch <- c.downloaded
	ch <- c.missing
	ch <- c.seriesSizeBytes
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SeriesCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	series, err := c.client.GetSeries(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "series")

	if err != nil {
		c.logger.Error("series collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "series")
		return
	}

	var monitoredCount int
	var totalEpisodes, downloadedEpisodes int
	var totalSeasons int
	var sizeOnDisk int64
	statusCounts := make(map[string]int)

	for _, s := range series {
		if s.Monitored {
			monitoredCount++
		}
		totalEpisodes += s.TotalEpisodeCount
		downloadedEpisodes += s.EpisodeFileCount
		totalSeasons += s.SeasonCount
		sizeOnDisk += s.SizeOnDisk
		if s.Status != "" {
			statusCounts[s.Status]++
		}
	}

	ch <- prometheus.MustNewConstMetric(c.seriesTotal, prometheus.GaugeValue, float64(len(series)))
	ch <- prometheus.MustNewConstMetric(c.monitored, prometheus.GaugeValue, float64(monitoredCount))
	for status, count := range statusCounts {
		ch <- prometheus.MustNewConstMetric(c.seriesByStatus, prometheus.GaugeValue, float64(count), status)
	}
	ch <- prometheus.MustNewConstMetric(c.seasonsTotal, prometheus.GaugeValue, float64(totalSeasons))
	ch <- prometheus.MustNewConstMetric(c.episodesTotal, prometheus.GaugeValue, float64(totalEpisodes))
	ch <- prometheus.MustNewConstMetric(c.downloaded, prometheus.GaugeValue, float64(downloadedEpisodes))
	ch <- prometheus.MustNewConstMetric(c.seriesSizeBytes, prometheus.GaugeValue, float64(sizeOnDisk))

	wanted, err := c.client.GetWantedMissing(ctx)
	if err != nil {
		c.logger.Error("wanted missing failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "series")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.missing, prometheus.GaugeValue, float64(wanted.TotalRecords))
	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "series")
}
