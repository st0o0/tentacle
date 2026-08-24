## Why

Tentacle currently only exports Jellyfin metrics. The same self-hosted media stack runs Sonarr, Radarr, Prowlarr, Seerr, and Audiobookshelf — each needing its own monitoring. A separate exporter per service (like the existing Conch project for Audiobookshelf) means multiple deployments, duplicated patterns, and fragmented dashboards. One unified exporter with a single `/metrics` endpoint simplifies deployment and gives a complete view of the media stack.

## What Changes

- **Restructure internal packages** from flat Jellyfin-only layout to per-service packages (`internal/client/<service>/`, `internal/collector/<service>/`)
- **Merge Conch** (Audiobookshelf exporter) into Tentacle as a new service module
- **Add shared Arr client** (`internal/client/arr/`) used by Sonarr, Radarr, and Prowlarr with service-specific extensions
- **Add Sonarr collectors** — system, health, series, episodes, queue, disk, calendar
- **Add Radarr collectors** — system, health, movies, queue, disk, calendar
- **Add Prowlarr collectors** — system, health, indexers, indexer stats
- **Add Seerr collectors** — system, requests, users
- **Extend config** to support per-service address/token pairs; a service is enabled when its env vars are set
- **Update Docker image and CI** to reflect the expanded scope

## Capabilities

### New Capabilities

- `arr-client`: Shared HTTP client for Sonarr/Radarr/Prowlarr APIs (X-Api-Key auth, /api/v3 base, shared types for system status, health, queue, disk)
- `audiobookshelf-client`: Audiobookshelf API client (ported from Conch)
- `seerr-client`: Seerr API client for media request tracking
- `collector-audiobookshelf`: Audiobookshelf collectors — system, libraries, users, sessions, backups
- `collector-sonarr`: Sonarr collectors — system, health, series, queue, disk, calendar
- `collector-radarr`: Radarr collectors — system, health, movies, queue, disk, calendar
- `collector-prowlarr`: Prowlarr collectors — system, health, indexers, indexer stats
- `collector-seerr`: Seerr collectors — system, requests, users

### Modified Capabilities

- `config`: Extend from single-service to multi-service configuration with per-service address/token env vars
- `jellyfin-client`: Move from `internal/jellyfin/` to `internal/client/jellyfin/` (structural change only)
- `metrics-server`: No functional change, but now serves metrics from all enabled services on the same registry

## Impact

- **Package structure**: `internal/jellyfin/` → `internal/client/jellyfin/`, `internal/collector/` split into per-service sub-packages
- **Config env vars**: All existing `TENTACLE_JELLYFIN_*` vars remain unchanged; new `TENTACLE_SONARR_*`, `TENTACLE_RADARR_*`, `TENTACLE_PROWLARR_*`, `TENTACLE_AUDIOBOOKSHELF_*`, `TENTACLE_SEERR_*` vars added
- **Docker**: Single image, same scratch base, more env vars documented
- **Conch project**: Will be superseded by Tentacle; Conch repo can be archived

## Non-goals

- Supporting multiple instances of the same service (e.g., two Sonarr servers)
- Config file support (YAML/TOML) — env vars only for now
- Aggregated cross-service metrics or dashboards (that's a Grafana concern)
- Write operations against any service API — Tentacle is read-only
- Lidarr, Readarr, Bazarr, or other arr-adjacent services (can be added later)
