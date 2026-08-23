## Context

Tentacle is a Jellyfin-only Prometheus exporter. A separate project (Conch) exports Audiobookshelf metrics using the identical architecture. The user runs a full self-hosted media stack (Jellyfin, Audiobookshelf, Sonarr, Radarr, Prowlarr, Seerr) and wants unified monitoring from a single exporter.

Both existing codebases use the same patterns: Go 1.26, `prometheus/client_golang`, `prometheus.Collector` interface, env var config with typed helpers, `slog` logging, `scratch` Docker image. The expansion should preserve these patterns.

## Goals / Non-Goals

**Goals:**
- Single binary exporting metrics for all 6 services on one `/metrics` endpoint
- Each service independently enabled via its env var pair (address + token)
- Shared Arr client for Sonarr/Radarr/Prowlarr to avoid code duplication
- Merge Conch's Audiobookshelf code into Tentacle
- Preserve all existing Jellyfin metrics and env var names

**Non-Goals:**
- Multi-instance support (multiple Sonarr servers, etc.)
- Config file support (YAML/TOML)
- Write operations against any service API
- Cross-service aggregated metrics
- Supporting Lidarr, Readarr, Bazarr

## Decisions

### Package structure: per-service sub-packages

Restructure from flat `internal/jellyfin/` and `internal/collector/` to:

```
internal/
  client/
    jellyfin/    ← moved from internal/jellyfin/
    audiobookshelf/
    arr/         ← shared client for Sonarr/Radarr/Prowlarr
    seerr/
  collector/
    collector.go ← shared scrape meta-metrics, namespace becomes a parameter
    jellyfin/    ← moved from internal/collector/*.go
    audiobookshelf/
    sonarr/
    radarr/
    prowlarr/
    seerr/
  config/
  metrics/
```

**Why over keeping flat:** A single `collector/` package with 30+ files and a single `namespace` constant doesn't scale. Per-service packages give each service its own namespace constant and import its own client type directly.

**Alternative considered:** Separate Go modules per service. Rejected because it adds build complexity for no benefit — this is one binary.

### Shared Arr client

Sonarr, Radarr, and Prowlarr are forks sharing the same API pattern (`/api/v3`, `X-Api-Key` header). A single `ArrClient` struct in `internal/client/arr/` provides:

- HTTP transport with `X-Api-Key` auth
- Shared endpoints: `GetSystemStatus()`, `GetHealth()`, `GetQueue()`, `GetRootFolders()`
- Shared types: `SystemStatus`, `HealthCheck`, `QueueRecord`, `RootFolder`

Service-specific methods live in the same package as typed extensions (e.g., `sonarr.go` adds `GetSeries()`, `GetEpisodes()`).

**Why over separate clients:** The three services share ~60% of their API surface. A shared base eliminates duplicated HTTP/auth/error-handling code and keeps shared types in sync.

**Why not an interface:** The Arr services genuinely share the same API — this isn't polymorphism, it's code reuse. One concrete struct instantiated three times with different URLs/tokens.

### Service activation by presence

A service is enabled when both `TENTACLE_<SERVICE>_ADDRESS` and `TENTACLE_<SERVICE>_TOKEN` are set. No explicit enable/disable flags.

**Why:** Self-documenting — if the vars aren't there, the service isn't configured. Adding explicit `TENTACLE_<SERVICE>_ENABLED` flags doubles the config surface for no benefit.

**Exception:** `TENTACLE_DISABLE_PLAYBACK` remains for the Jellyfin playback collector since it requires an optional plugin, not a separate service.

### Metric namespace per service

Each service gets its own Prometheus namespace prefix matching the service name:

| Service | Namespace | Example |
|---|---|---|
| Jellyfin | `jellyfin_` | `jellyfin_up` |
| Audiobookshelf | `audiobookshelf_` | `audiobookshelf_up` |
| Sonarr | `sonarr_` | `sonarr_series_total` |
| Radarr | `radarr_` | `radarr_movies_total` |
| Prowlarr | `prowlarr_` | `prowlarr_indexers_total` |
| Seerr | `seerr_` | `seerr_requests_total` |

Scrape meta-metrics also use the service namespace: `sonarr_scrape_duration_seconds{collector="system"}`.

**Why not `tentacle_`:** Users expect metric namespaces to identify the data source. `jellyfin_up` is immediately clear; `tentacle_jellyfin_up` adds noise.

### Collector registration in main.go

`main.go` creates a `ServiceConfig` per service from the parsed config. For each service where the config is present, it instantiates the client and registers all collectors. The pattern:

```
cfg := config.Load(os.Getenv)
reg := prometheus.NewRegistry()

if cfg.Jellyfin != nil {
    jc := jellyfin.NewClient(...)
    reg.MustRegister(jellyfinCollectors(jc, cfg.ScrapeTimeout, logger)...)
}
if cfg.Sonarr != nil {
    ac := arr.NewClient(cfg.Sonarr.Address, cfg.Sonarr.Token, httpClient)
    reg.MustRegister(sonarrCollectors(ac, cfg.ScrapeTimeout, logger)...)
}
// ... same for each service
```

### Shared collector helpers

The shared `collector.go` moves scrape meta-metric descriptors into a function that takes the namespace as a parameter:

```go
func NewScrapeDescs(namespace string) (duration, success *prometheus.Desc)
```

Each per-service collector package calls this with its own namespace.

## Risks / Trade-offs

- **Arr API stability** → Sonarr/Radarr/Prowlarr v3 API is stable and well-documented. Pin to `/api/v3` explicitly.
- **Scrape duration** → With 6 services, a single scrape hits many external APIs. Mitigated by per-collector timeouts (already implemented) and independent collector error handling (one service failing doesn't block others).
- **Binary size** → Adding 5 more clients increases binary size marginally. All clients are stdlib HTTP — no heavy dependencies added.
- **Conch deprecation** → After merge, Conch repo should be archived with a pointer to Tentacle. Existing Conch users need migration docs for env var changes (`CONCH_` → `TENTACLE_AUDIOBOOKSHELF_`).

## Open Questions

- None at this time. API surfaces for all services are well-documented and stable.
