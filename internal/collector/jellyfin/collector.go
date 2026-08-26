package jellyfin

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/collector"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

const namespace = "jellyfin"

type Option func(*options)

type options struct {
	playback bool
}

func WithPlayback(enabled bool) Option {
	return func(o *options) {
		o.playback = enabled
	}
}

func NewCollector(client *jellyfin.Client, timeout time.Duration, batchSize int, logger *slog.Logger, opts ...Option) *collector.ServiceCollector {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	subs := []collector.SubCollector{
		newSystemCollector(client, timeout, logger),
		newUsersCollector(client, timeout, logger),
		newSessionsCollector(client, timeout, logger),
		newLibraryCollector(client, timeout, batchSize, logger),
		newTasksCollector(client, timeout, logger),
		newActivityCollector(client, timeout, logger),
		newPluginsCollector(client, timeout, logger),
		newDevicesCollector(client, timeout, logger),
		newCountsCollector(client, timeout, logger),
	}

	if o.playback {
		subs = append(subs, newPlaybackCollector(client, timeout, logger))
	}

	return collector.NewServiceCollector(namespace, logger, subs...)
}

func parseJellyfinTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
