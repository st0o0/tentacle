package radarr

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/arr"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "radarr"

func NewCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		newSystemCollector(client, timeout, logger),
		newMoviesCollector(client, timeout, logger),
		newQueueCollector(client, timeout, logger),
		newDiskCollector(client, timeout, logger),
		newCalendarCollector(client, timeout, logger),
		newExtrasCollector(client, timeout, logger),
	)
}
