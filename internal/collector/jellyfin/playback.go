package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type playbackCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	playCount *prometheus.Desc
	watchTime *prometheus.Desc
}

func newPlaybackCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *playbackCollector {
	return &playbackCollector{
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

func (c *playbackCollector) Name() string { return "playback" }

func (c *playbackCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.playCount
	ch <- c.watchTime
}

func (c *playbackCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	activity, err := c.client.GetPlaybackActivity(ctx, 30)
	if err != nil {
		return err
	}

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
	return nil
}
