package audiobookshelf

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "audiobookshelf"

func NewCollector(client *audiobookshelf.Client, timeout time.Duration, cache collector.CacheIntervals, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		collector.NewCachedCollector(newSystemCollector(client, timeout, logger), cache.Cold, logger),
		collector.NewCachedCollector(newLibrariesCollector(client, timeout, logger), cache.Warm, logger),
		collector.NewCachedCollector(newUsersCollector(client, timeout, logger), cache.Warm, logger),
		newSessionsCollector(client, timeout, logger),
		collector.NewCachedCollector(newBackupsCollector(client, timeout, logger), cache.Cold, logger),
	)
}
