package audiobookshelf

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "audiobookshelf"

func NewCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		newSystemCollector(client, timeout, logger),
		newLibrariesCollector(client, timeout, logger),
		newUsersCollector(client, timeout, logger),
		newSessionsCollector(client, timeout, logger),
		newBackupsCollector(client, timeout, logger),
	)
}
