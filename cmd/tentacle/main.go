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
			addr := os.Getenv("TENTACLE_LISTEN_ADDRESS")
			if addr == "" {
				addr = ":9594"
			}
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	httpClient := &http.Client{Timeout: cfg.ScrapeTimeout}
	reg := prometheus.NewRegistry()

	if cfg.Jellyfin != nil {
		client := jellyfin.NewClient(cfg.Jellyfin.Address, cfg.Jellyfin.Token, httpClient)

		var opts []jellyfincollector.Option
		probeCtx, probeCancel := context.WithTimeout(context.Background(), cfg.ScrapeTimeout)
		_, probeErr := client.GetPlaybackActivity(probeCtx, 1)
		probeCancel()
		if probeErr == nil {
			opts = append(opts, jellyfincollector.WithPlayback(true))
			logger.Info("playback collector enabled (PlaybackReporting plugin detected)")
		} else {
			logger.Info("playback collector disabled (PlaybackReporting plugin not found)")
		}

		reg.MustRegister(jellyfincollector.NewCollector(client, cfg.ScrapeTimeout, logger, opts...))
		logger.Info("registered jellyfin collectors", "addr", cfg.Jellyfin.Address)
	}

	if cfg.Sonarr != nil {
		client := arr.NewClient(cfg.Sonarr.Address, cfg.Sonarr.Token, httpClient)
		reg.MustRegister(sonarrcollector.NewCollector(client, cfg.ScrapeTimeout, logger))
		logger.Info("registered sonarr collectors", "addr", cfg.Sonarr.Address)
	}

	if cfg.Radarr != nil {
		client := arr.NewClient(cfg.Radarr.Address, cfg.Radarr.Token, httpClient)
		reg.MustRegister(radarrcollector.NewCollector(client, cfg.ScrapeTimeout, logger))
		logger.Info("registered radarr collectors", "addr", cfg.Radarr.Address)
	}

	if cfg.Prowlarr != nil {
		client := arr.NewClient(cfg.Prowlarr.Address, cfg.Prowlarr.Token, httpClient)
		reg.MustRegister(prowlarrcollector.NewCollector(client, cfg.ScrapeTimeout, logger))
		logger.Info("registered prowlarr collectors", "addr", cfg.Prowlarr.Address)
	}

	if cfg.Audiobookshelf != nil {
		client := audiobookshelf.NewClient(cfg.Audiobookshelf.Address, cfg.Audiobookshelf.Token, httpClient)
		reg.MustRegister(abscollector.NewCollector(client, cfg.ScrapeTimeout, logger))
		logger.Info("registered audiobookshelf collectors", "addr", cfg.Audiobookshelf.Address)
	}

	if cfg.Seerr != nil {
		client := seerr.NewClient(cfg.Seerr.Address, cfg.Seerr.Token, httpClient)
		reg.MustRegister(seerrcollector.NewCollector(client, cfg.ScrapeTimeout, logger))
		logger.Info("registered seerr collectors", "addr", cfg.Seerr.Address)
	}

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
