package config

import (
	"log/slog"
	"testing"
	"time"
)

func envFrom(m map[string]string) func(string) string {
	return func(key string) string { return m[key] }
}

func TestLoad_NoServices(t *testing.T) {
	_, err := Load(envFrom(map[string]string{}))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got != "at least one service must be configured" {
		t.Fatalf("unexpected error: %s", got)
	}
}

func TestLoad_SingleService(t *testing.T) {
	cfg, err := Load(envFrom(map[string]string{
		"TENTACLE_JELLYFIN_ADDRESS": "http://jf:8096",
		"TENTACLE_JELLYFIN_TOKEN":   "jf-token",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Jellyfin == nil {
		t.Fatal("Jellyfin should be non-nil")
	}
	if cfg.Jellyfin.Address != "http://jf:8096" || cfg.Jellyfin.Token != "jf-token" {
		t.Fatalf("Jellyfin config mismatch: %+v", cfg.Jellyfin)
	}
	for name, svc := range map[string]*ServiceConfig{
		"Sonarr": cfg.Sonarr, "Radarr": cfg.Radarr, "Prowlarr": cfg.Prowlarr,
		"Audiobookshelf": cfg.Audiobookshelf, "Seerr": cfg.Seerr,
	} {
		if svc != nil {
			t.Fatalf("%s should be nil", name)
		}
	}
}

func TestLoad_MultipleServices(t *testing.T) {
	cfg, err := Load(envFrom(map[string]string{
		"TENTACLE_JELLYFIN_ADDRESS": "http://jf:8096",
		"TENTACLE_JELLYFIN_TOKEN":   "jf-token",
		"TENTACLE_SONARR_ADDRESS":   "http://sonarr:8989",
		"TENTACLE_SONARR_TOKEN":     "sonarr-token",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Jellyfin == nil {
		t.Fatal("Jellyfin should be non-nil")
	}
	if cfg.Sonarr == nil {
		t.Fatal("Sonarr should be non-nil")
	}
}

func TestLoad_PartialConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr string
	}{
		{
			name: "address without token",
			env: map[string]string{
				"TENTACLE_JELLYFIN_ADDRESS": "http://jf:8096",
			},
			wantErr: "TENTACLE_JELLYFIN_TOKEN is required when TENTACLE_JELLYFIN_ADDRESS is set",
		},
		{
			name: "token without address",
			env: map[string]string{
				"TENTACLE_JELLYFIN_TOKEN": "jf-token",
			},
			wantErr: "TENTACLE_JELLYFIN_ADDRESS is required when TENTACLE_JELLYFIN_TOKEN is set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(envFrom(tt.env))
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got := err.Error(); got != tt.wantErr {
				t.Fatalf("got error %q, want %q", got, tt.wantErr)
			}
		})
	}
}

func TestLoad_GlobalDefaults(t *testing.T) {
	cfg, err := Load(envFrom(map[string]string{
		"TENTACLE_JELLYFIN_ADDRESS": "http://jf:8096",
		"TENTACLE_JELLYFIN_TOKEN":   "jf-token",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenAddress != ":9594" {
		t.Fatalf("ListenAddress = %q, want %q", cfg.ListenAddress, ":9594")
	}
	if cfg.ScrapeTimeout != 10*time.Second {
		t.Fatalf("ScrapeTimeout = %v, want %v", cfg.ScrapeTimeout, 10*time.Second)
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Fatalf("LogLevel = %v, want %v", cfg.LogLevel, slog.LevelInfo)
	}
	if cfg.LogFormat != "json" {
		t.Fatalf("LogFormat = %q, want %q", cfg.LogFormat, "json")
	}
}

