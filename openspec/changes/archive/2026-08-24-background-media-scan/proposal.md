## Why

Media codec distribution, resolution breakdown, and container format statistics are high-value metrics for library management ("How far is my HEVC migration?", "What percentage is 4K?"). However, collecting this data requires iterating all library items with MediaSources — too expensive for every Prometheus scrape. This change introduces a background scanning infrastructure that runs expensive analyses on a configurable interval and serves cached results on the standard `/metrics` endpoint, then uses it for Jellyfin media analysis as the first implementation.

## What Changes

- **New infrastructure**: Cached collector pattern — a background goroutine that periodically runs expensive data collection and caches results. Standard Prometheus scrapes serve the cached data instantly.
- **New config**: `TENTACLE_SCAN_INTERVAL` environment variable (default `6h`, `0` to disable)
- **Jellyfin media analysis**: Per-library breakdown of video codecs, audio codecs, resolutions, container formats, and subtitle codecs
- **Jellyfin client**: New method to fetch items with MediaSources/MediaStreams fields
- **Scan meta-metrics**: Track when the last scan completed, how long it took, and how many items were analyzed

## Capabilities

### New Capabilities

- `cached-collector`: Infrastructure for background-scanned metrics with configurable interval and cache
- `collector-media-analysis`: Jellyfin media codec, resolution, and container analysis collector

### Modified Capabilities

- `config`: Add TENTACLE_SCAN_INTERVAL environment variable
- `jellyfin-client`: Add method to fetch items with MediaSources fields for stream analysis

## Impact

- **Architecture**: New cached collector pattern in internal/collector/ — reusable for future expensive metrics
- **Config**: One new environment variable (TENTACLE_SCAN_INTERVAL)
- **API load**: Significant Jellyfin API calls during scan (paginated item fetches with MediaSources), but only every N hours
- **Memory**: Cached metric results held in memory between scans — bounded by library size and label cardinality
- **Dependencies**: No new external dependencies

## Non-Goals

- No media analysis for non-Jellyfin services (Audiobookshelf doesn't expose codec info via API)
- No per-item metrics — only aggregated per-library counts
- No historical trend storage — Prometheus handles that
- No scan triggering via API — only interval-based
