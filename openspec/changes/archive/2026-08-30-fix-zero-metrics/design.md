## Context

tentacle v0.1.5 has two API-mapping bugs causing zero-value metrics:

1. **Sonarr**: The `/api/v3/series` endpoint returns episode/season counts and size inside a nested `statistics` object. The Go `Series` struct maps these fields at the top level, where Sonarr returns 0. Radarr works because its `/api/v3/movie` endpoint puts `sizeOnDisk` at the top level.

2. **Jellyfin**: The `/Items` endpoint's top-level `Size` field is not reliably populated. The actual file size lives in `MediaSources[].Size`, but the `MediaSource` struct has no `Size` field. Additionally, `librarySize()` queries all item types including containers (Series, Season) which always have Size=0.

## Goals / Non-Goals

**Goals:**
- Fix all zero-value Sonarr metrics: `episodes_total`, `episodes_downloaded_total`, `seasons_total`, `series_size_bytes`
- Fix zero-value Jellyfin metric: `library_size_bytes`

**Non-Goals:**
- `disk_total_bytes=0` (upstream API limitation in Docker/NFS)
- Jellyfin episode count discrepancy (by design: global dedup vs per-library)
- Adding new metrics

## Decisions

### D1: Sonarr — Add Statistics sub-struct

Map `statistics` as a nested struct on `Series`. Read counts and size from `s.Statistics.*` in the collector.

**Alternative considered:** Separate API call to `/api/v3/episode` per series — rejected because it would generate N+1 requests for 93+ series, while the data is already in the series response.

### D2: Jellyfin — Sum MediaSources sizes with leaf-item filter

Two changes:
1. Add `Size int64` to `MediaSource` struct.
2. Change `librarySize()` to request `Fields=MediaSources` instead of `Fields=Size`, and sum `MediaSources[].Size` per item. Also add `IncludeItemTypes` filter for leaf types (Movie, Episode, Audio, MusicVideo, Book) to skip containers that have no files.

**Alternative considered:** Keep using top-level `Item.Size` but filter to leaf types only — rejected because `Item.Size` is unreliable even for leaf types in some Jellyfin versions. `MediaSources[].Size` is the canonical source.

**Alternative considered:** Query `MediaSources` field in the existing per-type item count calls — rejected because those calls use `Limit=0` (count-only) and adding size there would require fetching all items for every type redundantly.

## Risks / Trade-offs

- **[Sonarr API version]** Older Sonarr v3 versions might populate top-level fields instead of `statistics`. → Mitigation: Keep both paths — if `Statistics` is zero, fall back to top-level fields. However, testing shows even Sonarr v3 uses `statistics`, so this is low risk. We'll map both and prefer `Statistics` when non-zero.

- **[Jellyfin MediaSources payload size]** Requesting `Fields=MediaSources` returns more data per item than `Fields=Size`. → Mitigation: The batch pagination (default 500) already handles large libraries. MediaSources adds ~100-200 bytes per item — negligible for batch sizes of 500.

- **[Items without MediaSources]** Some items (e.g., live TV recordings, external streams) may have empty MediaSources. → Mitigation: Sum defaults to 0 for those items, same behavior as today. No regression.
