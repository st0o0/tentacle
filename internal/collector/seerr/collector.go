package seerr

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/seerr"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "seerr"

func NewCollector(client *seerr.Client, timeout time.Duration, cache collector.CacheIntervals, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		collector.NewCachedCollector(newSystemCollector(client, timeout, logger), cache.Cold, logger),
		collector.NewCachedCollector(newRequestsCollector(client, timeout, logger), cache.Warm, logger),
		collector.NewCachedCollector(newUsersCollector(client, timeout, logger), cache.Cold, logger),
		collector.NewCachedCollector(newIssuesCollector(client, timeout, logger), cache.Warm, logger),
	)
}
