## Why

Several metrics in tentacle v0.1.5 report 0 despite the upstream APIs returning valid data. Sonarr's `/api/v3/series` nests episode counts and size under a `statistics` sub-object, but the Go struct only maps top-level fields (which are 0). Jellyfin's `/Items` endpoint populates size in `MediaSources[].Size`, not `Item.Size`, but the `MediaSource` struct lacks a `Size` field. These are the only two metrics families affected and both are API-mapping errors — the data is there, just read from the wrong location.

## What Changes

- **Sonarr series struct**: Add a `Statistics` sub-struct to `Series` in `sonarr_types.go` mapping `statistics.seasonCount`, `statistics.totalEpisodeCount`, `statistics.episodeFileCount`, and `statistics.sizeOnDisk`.
- **Sonarr series collector**: Read episode/season counts and size from `s.Statistics.*` instead of top-level fields in `series.go`.
- **Jellyfin MediaSource struct**: Add `Size int64` to the `MediaSource` struct in `types.go`.
- **Jellyfin library collector**: Compute library size from `MediaSources[0].Size` (or sum of all MediaSources) instead of top-level `Item.Size` in `library.go`.

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `collector-sonarr`: Series collector reads statistics from nested API field instead of top-level.
- `collector-library`: Library size calculation uses `MediaSources[].Size` instead of `Item.Size`.
- `arr-client`: Sonarr `Series` type gains a `Statistics` sub-struct.
- `jellyfin-client`: `MediaSource` type gains a `Size` field.

## Impact

- `internal/client/arr/sonarr_types.go` — struct change
- `internal/collector/sonarr/series.go` — field access change
- `internal/client/jellyfin/types.go` — struct change
- `internal/collector/jellyfin/library.go` — size calculation logic

No API changes, no config changes, no new dependencies.

## Non-goals

- Fixing `disk_total_bytes=0` — this is an upstream Sonarr/Radarr API limitation in Docker/NFS environments, not an exporter bug.
- Resolving the Jellyfin episode count discrepancy (8473 vs 10572) — this is by design: `items_total` deduplicates, `library_items_total` counts per-library.
- Adding new metrics or collectors.
