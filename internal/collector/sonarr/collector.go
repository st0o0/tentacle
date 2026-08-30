package sonarr

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/arr"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "sonarr"

func NewCollector(client *arr.Client, timeout time.Duration, cache collector.CacheIntervals, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		collector.NewCachedCollector(newSystemCollector(client, timeout, logger), cache.Cold, logger),
		collector.NewCachedCollector(newSeriesCollector(client, timeout, logger), cache.Warm, logger),
		newQueueCollector(client, timeout, logger),
		collector.NewCachedCollector(newDiskCollector(client, timeout, logger), cache.Cold, logger),
		collector.NewCachedCollector(newCalendarCollector(client, timeout, logger), cache.Warm, logger),
		collector.NewCachedCollector(newExtrasCollector(client, timeout, logger), cache.Cold, logger),
	)
}
