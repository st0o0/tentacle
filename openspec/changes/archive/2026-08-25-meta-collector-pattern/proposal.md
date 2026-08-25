## Why

Registering multiple collectors per service on the same `prometheus.Registry` panics because each collector's `Describe()` sends the shared `scrape_success` and `scrape_duration_seconds` descriptors — the registry treats these as duplicates from different collectors. This crash-loops the container on startup.

## What Changes

- Introduce a `SubCollector` interface with a single `Update(ch chan<- prometheus.Metric) error` method, replacing direct `prometheus.Collector` implementation in sub-collectors
- Add a per-service wrapper collector (e.g., `JellyfinCollector`) that implements `prometheus.Collector`, owns the scrape descriptors, orchestrates sub-collectors, and emits scrape timing/success metrics centrally
- Convert all existing sub-collectors (system, users, sessions, etc.) from `prometheus.Collector` to `SubCollector`
- Simplify `main.go` registration from N calls per service to one call per service
- Handle optional sub-collectors (e.g., jellyfin playback) via conditional registration on the wrapper

## Non-goals

- Changing metric names, labels, or semantics — all existing metrics remain identical
- Adding concurrency to sub-collector scrapes (can be done later; sequential is fine for now)
- Modifying the `metrics.ListenAndServe` server or health endpoint

## Capabilities

### New Capabilities
- `collector-orchestration`: Per-service wrapper collector that owns scrape descriptors, iterates sub-collectors, and emits scrape instrumentation centrally

### Modified Capabilities
- `metrics-server`: The "Scrape instrumentation" requirement moves from individual collectors to the service wrapper — same metrics, different emitter

## Impact

- `internal/collector/collector.go` — new `SubCollector` interface and wrapper type
- `internal/collector/jellyfin/*.go` — all collectors converted to `SubCollector`
- `internal/collector/sonarr/*.go` — same conversion
- `internal/collector/radarr/*.go` — same conversion
- `internal/collector/prowlarr/*.go` — same conversion
- `internal/collector/audiobookshelf/*.go` — same conversion
- `internal/collector/seerr/*.go` — same conversion
- `cmd/tentacle/main.go` — simplified registration, playback probe logic moves into jellyfin wrapper
- Existing tests need updating to use the wrapper or test `Update()` directly
