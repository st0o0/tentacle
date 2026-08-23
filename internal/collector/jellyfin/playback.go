package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type PlaybackCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	playCount *prometheus.Desc
	watchTime *prometheus.Desc
}

func NewPlaybackCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *PlaybackCollector {
	return &PlaybackCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		playCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "playback", "play_count"),
			"Total play count per user (last 30 days, requires PlaybackReporting plugin).",
			[]string{"user"}, nil,
		),
		watchTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "playback", "watch_time_seconds"),
			"Total watch time in seconds per user (last 30 days, requires PlaybackReporting plugin).",
			[]string{"user"}, nil,
		),
	}
}

func (c *PlaybackCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.playCount
	ch <- c.watchTime
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *PlaybackCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	activity, err := c.client.GetPlaybackActivity(ctx, 30)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "playback")

	if err != nil {
		c.logger.Debug("playback collector unavailable (requires PlaybackReporting plugin)", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "playback")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "playback")

	type userStats struct {
		playCount int
		watchTime float64
	}
	users := make(map[string]*userStats)

	for _, a := range activity {
		name := a.UserName
		if name == "" {
			continue
		}
		s, ok := users[name]
		if !ok {
			s = &userStats{}
			users[name] = s
		}
		s.playCount += a.PlayCount
		s.watchTime += a.WatchTime
	}

	for user, s := range users {
		ch <- prometheus.MustNewConstMetric(c.playCount, prometheus.GaugeValue, float64(s.playCount), user)
		ch <- prometheus.MustNewConstMetric(c.watchTime, prometheus.GaugeValue, s.watchTime, user)
	}
}
