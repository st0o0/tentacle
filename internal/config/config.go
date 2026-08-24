package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type ServiceConfig struct {
	Address string
	Token   string
}

type Config struct {
	Jellyfin       *ServiceConfig
	Sonarr         *ServiceConfig
	Radarr         *ServiceConfig
	Prowlarr       *ServiceConfig
	Audiobookshelf *ServiceConfig
	Seerr          *ServiceConfig
	ListenAddress string
	ScrapeTimeout time.Duration
	LogLevel      slog.Level
	LogFormat     string
}

func Load(getenv func(string) string) (Config, error) {
	r := &envReader{getenv: getenv}

	cfg := Config{
		Jellyfin:        r.service("TENTACLE_JELLYFIN_ADDRESS", "TENTACLE_JELLYFIN_TOKEN"),
		Sonarr:          r.service("TENTACLE_SONARR_ADDRESS", "TENTACLE_SONARR_TOKEN"),
		Radarr:          r.service("TENTACLE_RADARR_ADDRESS", "TENTACLE_RADARR_TOKEN"),
		Prowlarr:        r.service("TENTACLE_PROWLARR_ADDRESS", "TENTACLE_PROWLARR_TOKEN"),
		Audiobookshelf:  r.service("TENTACLE_AUDIOBOOKSHELF_ADDRESS", "TENTACLE_AUDIOBOOKSHELF_TOKEN"),
		Seerr:           r.service("TENTACLE_SEERR_ADDRESS", "TENTACLE_SEERR_TOKEN"),
		ListenAddress:   r.str("TENTACLE_LISTEN_ADDRESS", ":9594"),
		ScrapeTimeout:   r.duration("TENTACLE_SCRAPE_TIMEOUT", 10*time.Second),
		LogLevel:        r.logLevel("TENTACLE_LOG_LEVEL", slog.LevelInfo),
		LogFormat:       r.logFormat("TENTACLE_LOG_FORMAT", "json"),
	}

	if r.err == nil && cfg.Jellyfin == nil && cfg.Sonarr == nil && cfg.Radarr == nil && cfg.Prowlarr == nil && cfg.Audiobookshelf == nil && cfg.Seerr == nil {
		r.setErr(errors.New("at least one service must be configured"))
	}

	return cfg, r.err
}

type envReader struct {
	getenv func(string) string
	err    error
}

func (r *envReader) setErr(err error) {
	if r.err == nil {
		r.err = err
	}
}

func (r *envReader) str(key, defaultVal string) string {
	v := r.getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}

func (r *envReader) service(addressKey, tokenKey string) *ServiceConfig {
	addr := r.getenv(addressKey)
	token := r.getenv(tokenKey)

	if addr == "" && token == "" {
		return nil
	}

	if addr == "" {
		r.setErr(fmt.Errorf("%s is required when %s is set", addressKey, tokenKey))
		return nil
	}

	if token == "" {
		r.setErr(fmt.Errorf("%s is required when %s is set", tokenKey, addressKey))
		return nil
	}

	return &ServiceConfig{Address: addr, Token: token}
}

func (r *envReader) duration(key string, defaultVal time.Duration) time.Duration {
	v := r.getenv(key)
	if v == "" {
		return defaultVal
	}

	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		r.setErr(fmt.Errorf("%s: invalid duration %q", key, v))
		return defaultVal
	}
	return d
}

func (r *envReader) logLevel(key string, defaultVal slog.Level) slog.Level {
	v := r.getenv(key)
	if v == "" {
		return defaultVal
	}
	switch strings.ToLower(v) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		r.setErr(fmt.Errorf("%s: invalid log level %q (debug, info, warn, error)", key, v))
		return defaultVal
	}
}

func (r *envReader) logFormat(key, defaultVal string) string {
	v := r.getenv(key)
	if v == "" {
		return defaultVal
	}
	switch strings.ToLower(v) {
	case "json", "text":
		return strings.ToLower(v)
	default:
		r.setErr(fmt.Errorf("%s: invalid log format %q (json, text)", key, v))
		return defaultVal
	}
}

