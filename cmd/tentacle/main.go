package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/arr"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
	"github.com/st0o0/tentacle/internal/client/seerr"
	abscollector "github.com/st0o0/tentacle/internal/collector/audiobookshelf"
	jellyfincollector "github.com/st0o0/tentacle/internal/collector/jellyfin"
	prowlarrcollector "github.com/st0o0/tentacle/internal/collector/prowlarr"
	radarrcollector "github.com/st0o0/tentacle/internal/collector/radarr"
	seerrcollector "github.com/st0o0/tentacle/internal/collector/seerr"
	sonarrcollector "github.com/st0o0/tentacle/internal/collector/sonarr"
	"github.com/st0o0/tentacle/internal/config"
	"github.com/st0o0/tentacle/internal/metrics"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println(version)
			os.Exit(0)
		case "healthcheck":
			addr := ":9594"
			if len(os.Args) > 2 {
				addr = os.Args[2]
			}
			os.Exit(runHealthcheck(addr))
		}
	}

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	logger := newLogger(cfg)
	slog.SetDefault(logger)

	httpClient := &http.Client{Timeout: cfg.ScrapeTimeout}
	reg := prometheus.NewRegistry()

	if cfg.Jellyfin != nil {
		client := jellyfin.NewClient(cfg.Jellyfin.Address, cfg.Jellyfin.Token, httpClient)
		reg.MustRegister(
			jellyfincollector.NewSystemCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewUsersCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewSessionsCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewLibraryCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewTasksCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewActivityCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewPluginsCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewDevicesCollector(client, cfg.ScrapeTimeout, logger),
			jellyfincollector.NewCountsCollector(client, cfg.ScrapeTimeout, logger),
		)

		probeCtx, probeCancel := context.WithTimeout(context.Background(), cfg.ScrapeTimeout)
		_, probeErr := client.GetPlaybackActivity(probeCtx, 1)
		probeCancel()
		if probeErr == nil {
			reg.MustRegister(jellyfincollector.NewPlaybackCollector(client, cfg.ScrapeTimeout, logger))
			logger.Info("playback collector enabled (PlaybackReporting plugin detected)")
		} else {
			logger.Info("playback collector disabled (PlaybackReporting plugin not found)")
		}

		logger.Info("registered jellyfin collectors", "addr", cfg.Jellyfin.Address)
	}

	if cfg.Sonarr != nil {
		client := arr.NewClient(cfg.Sonarr.Address, cfg.Sonarr.Token, httpClient)
		reg.MustRegister(
			sonarrcollector.NewSystemCollector(client, cfg.ScrapeTimeout, logger),
			sonarrcollector.NewSeriesCollector(client, cfg.ScrapeTimeout, logger),
			sonarrcollector.NewQueueCollector(client, cfg.ScrapeTimeout, logger),
			sonarrcollector.NewDiskCollector(client, cfg.ScrapeTimeout, logger),
			sonarrcollector.NewCalendarCollector(client, cfg.ScrapeTimeout, logger),
		)
		logger.Info("registered sonarr collectors", "addr", cfg.Sonarr.Address)
	}

	if cfg.Radarr != nil {
		client := arr.NewClient(cfg.Radarr.Address, cfg.Radarr.Token, httpClient)
		reg.MustRegister(
			radarrcollector.NewSystemCollector(client, cfg.ScrapeTimeout, logger),
			radarrcollector.NewMoviesCollector(client, cfg.ScrapeTimeout, logger),
			radarrcollector.NewQueueCollector(client, cfg.ScrapeTimeout, logger),
			radarrcollector.NewDiskCollector(client, cfg.ScrapeTimeout, logger),
			radarrcollector.NewCalendarCollector(client, cfg.ScrapeTimeout, logger),
		)
		logger.Info("registered radarr collectors", "addr", cfg.Radarr.Address)
	}

	if cfg.Prowlarr != nil {
		client := arr.NewClient(cfg.Prowlarr.Address, cfg.Prowlarr.Token, httpClient)
		reg.MustRegister(
			prowlarrcollector.NewSystemCollector(client, cfg.ScrapeTimeout, logger),
			prowlarrcollector.NewIndexersCollector(client, cfg.ScrapeTimeout, logger),
		)
		logger.Info("registered prowlarr collectors", "addr", cfg.Prowlarr.Address)
	}

	if cfg.Audiobookshelf != nil {
		client := audiobookshelf.NewClient(cfg.Audiobookshelf.Address, cfg.Audiobookshelf.Token, httpClient)
		reg.MustRegister(
			abscollector.NewSystemCollector(client, cfg.ScrapeTimeout, logger),
			abscollector.NewLibrariesCollector(client, cfg.ScrapeTimeout, logger),
			abscollector.NewUsersCollector(client, cfg.ScrapeTimeout, logger),
			abscollector.NewSessionsCollector(client, cfg.ScrapeTimeout, logger),
			abscollector.NewBackupsCollector(client, cfg.ScrapeTimeout, logger),
		)
		logger.Info("registered audiobookshelf collectors", "addr", cfg.Audiobookshelf.Address)
	}

	if cfg.Seerr != nil {
		client := seerr.NewClient(cfg.Seerr.Address, cfg.Seerr.Token, httpClient)
		reg.MustRegister(
			seerrcollector.NewSystemCollector(client, cfg.ScrapeTimeout, logger),
			seerrcollector.NewRequestsCollector(client, cfg.ScrapeTimeout, logger),
			seerrcollector.NewUsersCollector(client, cfg.ScrapeTimeout, logger),
		)
		logger.Info("registered seerr collectors", "addr", cfg.Seerr.Address)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger.Info("starting tentacle", "version", version, "addr", cfg.ListenAddress)

	if err := metrics.ListenAndServe(ctx, cfg.ListenAddress, reg, version, logger); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}

func newLogger(cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	var handler slog.Handler
	switch cfg.LogFormat {
	case "text":
		handler = slog.NewTextHandler(os.Stderr, opts)
	default:
		handler = slog.NewJSONHandler(os.Stderr, opts)
	}
	return slog.New(handler)
}
