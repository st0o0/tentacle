package prowlarr

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/arr"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "prowlarr"

func NewCollector(client *arr.Client, timeout time.Duration, cache collector.CacheIntervals, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		collector.NewCachedCollector(newSystemCollector(client, timeout, logger), cache.Cold, logger),
		collector.NewCachedCollector(newIndexersCollector(client, timeout, logger), cache.Warm, logger),
		collector.NewCachedCollector(newAppsCollector(client, timeout, logger), cache.Warm, logger),
	)
}
