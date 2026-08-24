## Context

Prometheus scrapes `/metrics` every 15-60 seconds. Most collectors make 1-3 lightweight API calls per scrape. Media analysis requires paginating through all library items with `Fields=MediaSources` — potentially hundreds of API calls for large libraries. This is incompatible with the standard scrape-driven collector model.

The solution is a cached collector pattern: a background goroutine runs the expensive analysis on a configurable interval, stores results in memory, and the Prometheus collect method simply returns the cached metrics.

## Goals / Non-Goals

**Goals:**
- Create a reusable cached collector pattern for any future expensive metric
- Implement Jellyfin media analysis (codecs, resolution, containers) as the first use case
- Keep it transparent to Prometheus — single `/metrics` endpoint, no extra scrape config
- Make it configurable and disableable

**Non-Goals:**
- No on-demand scan triggering (no API endpoint to force a scan)
- No persistent cache (results are lost on restart, rebuilt on first scan)
- No incremental scanning (full scan each interval)

## Decisions

### Decision 1: Background goroutine with in-memory cache

**Choice**: A `CachedCollector` struct wraps a scan function and holds cached `[]prometheus.Metric` behind a `sync.RWMutex`. A background goroutine calls the scan function on a ticker, replaces the cache on success, and the `Collect()` method reads from cache under RLock.

**Why**: Simple, no external dependencies, fits the existing collector pattern. The cache is just a slice of pre-built metrics — no complex data structures.

**Alternative**: Separate `/metrics/slow` endpoint with its own Prometheus scrape job — rejected because it requires users to configure a second scrape job in Prometheus, which is unintuitive.

**Alternative**: Time-based cache in the normal collector (check elapsed time on each scrape, refresh if stale) — rejected because the scan takes seconds/minutes, which would cause scrape timeouts for the unlucky scrape that triggers the refresh.

### Decision 2: Scan interval configuration

**Choice**: `TENTACLE_SCAN_INTERVAL` with default `6h`. Value `0` disables the scan entirely. Accepts Go duration format.

**Why**: 6 hours is a reasonable default — library composition changes slowly. Disabling is important for users who don't want the API load. Using the existing duration parsing from config is consistent.

### Decision 3: Items pagination with MediaSources

**Choice**: Fetch items per library using `GetItems()` with `Fields=MediaSources` in batches of 500. Parse `MediaSources[].MediaStreams[]` to extract video codec, audio codec, resolution, container, and subtitle codec.

**Why**: `MediaSources` embeds all stream information in the item response — no separate API call per item needed. Batch size of 500 matches the existing library size calculation pattern.

**Resolution bucketing**: Map width to standard labels:
- width >= 3840 → "4k"
- width >= 1920 → "1080p"
- width >= 1280 → "720p"
- width >= 720 → "480p"
- else → "sd"

### Decision 4: Metric structure

**Choice**: Separate gauge vectors per dimension, all with `{library}` label:

```
jellyfin_library_video_codec_total{library, codec}
jellyfin_library_audio_codec_total{library, codec}
jellyfin_library_resolution_total{library, resolution}
jellyfin_library_container_total{library, container}
jellyfin_library_subtitle_codec_total{library, codec}
```

**Why**: Separate vectors allow independent queries ("show me resolution distribution" without filtering out codec data). The `library` label ties back to existing library metrics.

### Decision 5: Scan meta-metrics

**Choice**: Three meta-metrics emitted by the cached collector infrastructure:

```
jellyfin_scan_last_completed_timestamp_seconds   — when the last scan finished
jellyfin_scan_duration_seconds                   — how long the last scan took
jellyfin_scan_items_analyzed_total               — how many items were analyzed
```

**Why**: Essential for monitoring the scanner itself — "is the scan running?", "is it getting slower?", alerting on "scan hasn't completed in 24h".

### Decision 6: First scan timing

**Choice**: Start the first scan immediately on startup (in the background goroutine), then repeat on the configured interval. Before the first scan completes, the collector returns no media analysis metrics.

**Why**: Users get metrics as soon as possible after startup. The "no metrics until first scan" behavior is clear and expected — Prometheus handles missing metrics gracefully.

### Decision 7: Error handling during scan

**Choice**: If the scan fails (API errors), keep the previous cached results and log the error. Emit a `jellyfin_scan_success` gauge (1 on success, 0 on failure).

**Why**: Transient API failures shouldn't wipe existing metrics. The success gauge enables alerting on scan failures.

## Risks / Trade-offs

- **[API load during scan]** → A library with 5000 items needs ~10 paginated calls with full MediaSources data. Mitigation: Runs only every 6h by default; configurable; disableable.
- **[Memory usage]** → Cached metrics held in memory. For a library with 20 unique codecs × 5 libraries, this is ~100 metric objects — negligible.
- **[Stale data]** → Metrics can be up to `SCAN_INTERVAL` old. Mitigation: Acceptable for library composition which changes slowly. Meta-metrics show data age.
- **[First scan delay]** → No media metrics until first scan completes (could take 30s+ for large libraries). Mitigation: Documented behavior, consistent with expected startup sequence.

## Open Questions

_None._
