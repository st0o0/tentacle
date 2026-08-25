package seerr

import (
	"log/slog"
	"time"

	"github.com/st0o0/tentacle/internal/client/seerr"
	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "seerr"

func NewCollector(client *seerr.Client, timeout time.Duration, logger *slog.Logger) *collector.ServiceCollector {
	return collector.NewServiceCollector(namespace, logger,
		newSystemCollector(client, timeout, logger),
		newRequestsCollector(client, timeout, logger),
		newUsersCollector(client, timeout, logger),
		newIssuesCollector(client, timeout, logger),
	)
}
