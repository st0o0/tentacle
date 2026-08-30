package jellyfin

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/jellyfin"
	"github.com/st0o0/tentacle/internal/collector"
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

func NewCollector(client *jellyfin.Client, timeout time.Duration, batchSize int, cache collector.CacheIntervals, logger *slog.Logger, opts ...Option) *collector.ServiceCollector {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	// cold: library (expensive size calculation), system, plugins, devices
	coldTimeout := 2 * time.Minute
	library := newLibraryCollector(client, coldTimeout, batchSize, logger)
	system := newSystemCollector(client, timeout, logger)
	plugins := newPluginsCollector(client, timeout, logger)
	devices := newDevicesCollector(client, timeout, logger)

	// warm: users, activity, counts, playback
	users := newUsersCollector(client, timeout, logger)
	activity := newActivityCollector(client, timeout, logger)
	counts := newCountsCollector(client, timeout, logger)

	// live: sessions, tasks
	sessions := newSessionsCollector(client, timeout, logger)
	tasks := newTasksCollector(client, timeout, logger)

	subs := []collector.SubCollector{
		collector.NewCachedCollector(system, cache.Cold, logger),
		sessions,
		tasks,
		collector.NewCachedCollector(library, cache.Cold, logger),
		collector.NewCachedCollector(users, cache.Warm, logger),
		collector.NewCachedCollector(activity, cache.Warm, logger),
		collector.NewCachedCollector(plugins, cache.Cold, logger),
		collector.NewCachedCollector(devices, cache.Cold, logger),
		collector.NewCachedCollector(counts, cache.Warm, logger),
	}

	if o.playback {
		playback := newPlaybackCollector(client, timeout, logger)
		subs = append(subs, collector.NewCachedCollector(playback, cache.Warm, logger))
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
