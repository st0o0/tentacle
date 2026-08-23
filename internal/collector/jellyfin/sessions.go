package jellyfin

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type SessionsCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	activeTotal    *prometheus.Desc
	streamingTotal *prometheus.Desc
	transcoding    *prometheus.Desc
	bitrate        *prometheus.Desc
	progress       *prometheus.Desc
	paused         *prometheus.Desc
	sessionInfo    *prometheus.Desc
	transcodeInfo  *prometheus.Desc
	transcodeFps   *prometheus.Desc
	transcodeCompl *prometheus.Desc
	bandwidthTotal *prometheus.Desc
	directPlay     *prometheus.Desc
	transcodeCount *prometheus.Desc
}

func NewSessionsCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *SessionsCollector {
	return &SessionsCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		activeTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sessions", "active_total"),
			"Total number of active sessions.",
			nil, nil,
		),
		streamingTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sessions", "streaming_total"),
			"Number of sessions currently streaming.",
			[]string{"user", "play_method"}, nil,
		),
		transcoding: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "transcoding"),
			"Whether a session is transcoding.",
			[]string{"user", "media_type", "hw_acceleration"}, nil,
		),
		bitrate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "bitrate_bps"),
			"Current bitrate of a streaming session in bits per second.",
			[]string{"user"}, nil,
		),
		progress: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "progress_ratio"),
			"Playback progress as a ratio from 0 to 1.",
			[]string{"user", "media_title"}, nil,
		),
		paused: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "paused"),
			"Whether a session is paused.",
			[]string{"user", "media_title"}, nil,
		),
		sessionInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "info"),
			"Active session information.",
			[]string{"user", "device", "client", "client_version"}, nil,
		),
		transcodeInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "transcode_info"),
			"Transcode details for an active session.",
			[]string{"user", "video_codec", "audio_codec", "container", "resolution", "transcode_reason"}, nil,
		),
		transcodeFps: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "transcode_framerate"),
			"Current transcode framerate.",
			[]string{"user"}, nil,
		),
		transcodeCompl: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "transcode_completion_ratio"),
			"How much of the file has been pre-transcoded (0 to 1).",
			[]string{"user"}, nil,
		),
		bandwidthTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sessions", "bandwidth_total_bps"),
			"Total bandwidth of all active streaming sessions in bits per second.",
			nil, nil,
		),
		directPlay: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sessions", "direct_play_total"),
			"Number of sessions using direct play.",
			nil, nil,
		),
		transcodeCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "sessions", "transcode_total"),
			"Number of sessions using transcoding.",
			nil, nil,
		),
	}
}

func (c *SessionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.activeTotal
	ch <- c.streamingTotal
	ch <- c.transcoding
	ch <- c.bitrate
	ch <- c.progress
	ch <- c.paused
	ch <- c.sessionInfo
	ch <- c.transcodeInfo
	ch <- c.transcodeFps
	ch <- c.transcodeCompl
	ch <- c.bandwidthTotal
	ch <- c.directPlay
	ch <- c.transcodeCount
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *SessionsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	sessions, err := c.client.GetSessions(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "sessions")

	if err != nil {
		c.logger.Error("sessions collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "sessions")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "sessions")
	ch <- prometheus.MustNewConstMetric(c.activeTotal, prometheus.GaugeValue, float64(len(sessions)))

	type streamKey struct{ user, playMethod string }
	streamCounts := make(map[streamKey]float64)
	bitrateSum := make(map[string]float64)
	transcodeFpsSum := make(map[string]float64)
	transcodeComplSum := make(map[string]float64)

	var totalBandwidth float64
	var directPlayCount, transcodeCountVal float64

	for _, s := range sessions {
		if s.NowPlayingItem == nil {
			continue
		}

		playMethod := "Unknown"
		if s.PlayState != nil && s.PlayState.PlayMethod != "" {
			playMethod = s.PlayState.PlayMethod
		}

		user := s.UserName
		if user == "" {
			user = "unknown"
		}

		streamCounts[streamKey{user, playMethod}]++

		switch playMethod {
		case "DirectPlay", "DirectStream":
			directPlayCount++
		case "Transcode":
			transcodeCountVal++
		}

		ch <- prometheus.MustNewConstMetric(c.sessionInfo, prometheus.GaugeValue, 1, user, s.DeviceName, s.Client, s.ApplicationVersion)

		mediaType := s.NowPlayingItem.Type
		hwAccel := "false"
		isTranscoding := 0.0

		if s.TranscodingInfo != nil {
			if !s.TranscodingInfo.IsVideoDirect {
				isTranscoding = 1.0
			}
			if s.TranscodingInfo.HardwareAccelerationType != "" {
				hwAccel = "true"
			}
			bitrateSum[user] += float64(s.TranscodingInfo.Bitrate)
			totalBandwidth += float64(s.TranscodingInfo.Bitrate)
			transcodeFpsSum[user] += s.TranscodingInfo.Framerate
			transcodeComplSum[user] = s.TranscodingInfo.CompletionPercentage / 100.0

			resolution := fmt.Sprintf("%dx%d", s.TranscodingInfo.Width, s.TranscodingInfo.Height)
			reason := strings.Join(s.TranscodingInfo.TranscodeReasons, ",")
			ch <- prometheus.MustNewConstMetric(c.transcodeInfo, prometheus.GaugeValue, 1,
				user, s.TranscodingInfo.VideoCodec, s.TranscodingInfo.AudioCodec,
				s.TranscodingInfo.Container, resolution, reason)
		}

		ch <- prometheus.MustNewConstMetric(c.transcoding, prometheus.GaugeValue, isTranscoding, user, mediaType, hwAccel)

		if s.PlayState != nil {
			pausedVal := 0.0
			if s.PlayState.IsPaused {
				pausedVal = 1.0
			}
			ch <- prometheus.MustNewConstMetric(c.paused, prometheus.GaugeValue, pausedVal, user, s.NowPlayingItem.Name)

			if s.NowPlayingItem.RunTimeTicks > 0 {
				ratio := float64(s.PlayState.PositionTicks) / float64(s.NowPlayingItem.RunTimeTicks)
				ch <- prometheus.MustNewConstMetric(c.progress, prometheus.GaugeValue, ratio, user, s.NowPlayingItem.Name)
			}
		}
	}

	for key, count := range streamCounts {
		ch <- prometheus.MustNewConstMetric(c.streamingTotal, prometheus.GaugeValue, count, key.user, key.playMethod)
	}
	for user, bps := range bitrateSum {
		ch <- prometheus.MustNewConstMetric(c.bitrate, prometheus.GaugeValue, bps, user)
	}
	for user, fps := range transcodeFpsSum {
		ch <- prometheus.MustNewConstMetric(c.transcodeFps, prometheus.GaugeValue, fps, user)
	}
	for user, compl := range transcodeComplSum {
		ch <- prometheus.MustNewConstMetric(c.transcodeCompl, prometheus.GaugeValue, compl, user)
	}

	ch <- prometheus.MustNewConstMetric(c.bandwidthTotal, prometheus.GaugeValue, totalBandwidth)
	ch <- prometheus.MustNewConstMetric(c.directPlay, prometheus.GaugeValue, directPlayCount)
	ch <- prometheus.MustNewConstMetric(c.transcodeCount, prometheus.GaugeValue, transcodeCountVal)
}
