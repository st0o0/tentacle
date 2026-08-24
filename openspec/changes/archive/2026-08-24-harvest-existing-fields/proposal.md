## Why

Client methods across all services already fetch response fields that collectors silently discard — StartTime, Status, Protocol, CollectionType, IsActive, and more. These represent ~14 free Prometheus metrics with zero additional API load. Extracting them closes monitoring gaps (uptime tracking, status breakdowns, version consistency) without any config or infrastructure changes.

## What Changes

- **Arr services**: Expose `StartTime` as uptime gauge; break down series/movies by status; count Sonarr seasons; segment Prowlarr indexers by protocol
- **Jellyfin**: Add `collection_type` label to library metrics; emit missing item count types (Program, Item); add `app_version` label to device metrics; break down activity log by type
- **Audiobookshelf**: Call existing unused `GetStatus()` for system_info metric; count active vs inactive users; expose library last-update timestamps; call existing unused `GetUserListeningStats()` for per-user listening time

## Capabilities

### New Capabilities

_None — all changes extend existing collectors._

### Modified Capabilities

- `collector-sonarr`: Add start_time_seconds, series_by_status, seasons_total metrics from existing response fields
- `collector-radarr`: Add start_time_seconds, movies_by_status metrics from existing response fields
- `collector-prowlarr`: Add start_time_seconds, indexers_by_protocol metrics from existing response fields
- `collector-library`: Add collection_type label from VirtualFolder.CollectionType
- `collector-counts`: Emit ProgramCount and ItemCount (currently omitted)
- `collector-devices`: Add app_version label from DeviceInfo.AppVersion
- `collector-activity`: Add entries_by_type breakdown from ActivityLog.Type
- `collector-audiobookshelf`: Add system_info, users_active, library_last_update, user_listening_time metrics

## Impact

- **Collectors**: 8 existing collector specs modified with new metric requirements
- **Clients**: No changes — all data already returned by existing methods
- **Config**: No changes
- **Labels**: New labels on existing metrics (collection_type, app_version) are additive and non-breaking
- **API load**: Zero additional API calls for Arr/Jellyfin; ABS adds GetStatus() (1 call) and GetUserListeningStats() (N calls, one per user — typically <10)

## Non-Goals

- No new API endpoints or client methods (except invoking existing unused ones)
- No new collector files — only extend existing ones
- No config changes or new environment variables
- No changes to scrape interval or caching behavior
