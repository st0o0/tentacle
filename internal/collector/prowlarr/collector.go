package prowlarr

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/arr"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "prowlarr"

func NewCollector(client *arr.Client, timeout time.Duration, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		newSystemCollector(client, timeout, logger),
		newIndexersCollector(client, timeout, logger),
		newAppsCollector(client, timeout, logger),
	)
}
