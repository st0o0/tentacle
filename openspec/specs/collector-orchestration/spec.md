## Purpose

Per-service collector orchestration using the meta-collector pattern. A single `ServiceCollector` per service owns scrape descriptors and iterates sub-collectors, avoiding duplicate descriptor panics on the Prometheus registry.

## Requirements

### Requirement: SubCollector interface
Each domain collector SHALL implement a `SubCollector` interface with `Name() string`, `Describe(ch chan<- *prometheus.Desc)`, and `Update(ch chan<- prometheus.Metric) error`. Sub-collectors SHALL NOT emit scrape instrumentation metrics.

#### Scenario: Sub-collector emits only domain metrics
- **WHEN** a sub-collector's `Update` is called
- **THEN** it emits only its domain-specific metrics (e.g., `jellyfin_system_info`) and returns nil on success or an error on failure

#### Scenario: Sub-collector describes only domain descriptors
- **WHEN** a sub-collector's `Describe` is called
- **THEN** it sends only its domain-specific `*prometheus.Desc` values, not scrape duration or success descriptors

### Requirement: ServiceCollector wrapper
A `ServiceCollector` struct SHALL implement `prometheus.Collector`, own the scrape descriptors for its namespace, and orchestrate a list of `SubCollector` instances.

#### Scenario: Single registry registration per service
- **WHEN** a service (e.g., jellyfin) is configured
- **THEN** exactly one `ServiceCollector` is registered with the `prometheus.Registry`

#### Scenario: Describe aggregates all descriptors
- **WHEN** `ServiceCollector.Describe` is called
- **THEN** it sends the scrape duration and success descriptors once, plus all domain descriptors from each sub-collector

#### Scenario: Collect iterates sub-collectors concurrently
- **WHEN** `ServiceCollector.Collect` is called
- **THEN** it calls each sub-collector's `Update` concurrently, measures its duration, and emits `<namespace>_scrape_duration_seconds{collector="<name>"}` and `<namespace>_scrape_success{collector="<name>"}` per sub-collector

#### Scenario: Sub-collector failure does not abort others
- **WHEN** one sub-collector's `Update` returns an error
- **THEN** the error is logged, `scrape_success` is 0 for that collector, and remaining sub-collectors are still executed

### Requirement: CachedCollector wrapper
A `CachedCollector` SHALL wrap any `SubCollector` and refresh its metrics in the background on a configurable interval. During a Prometheus scrape it SHALL serve the most recently cached values instead of calling the upstream API.

#### Scenario: Background refresh on interval
- **WHEN** a `CachedCollector` is started with a 10-minute interval
- **THEN** it performs an immediate first refresh and then refreshes every 10 minutes in the background

#### Scenario: Scrape serves cached metrics
- **WHEN** `CachedCollector.Update` is called during a Prometheus scrape
- **THEN** it returns the cached metrics from the last background refresh without calling the inner collector

#### Scenario: Empty cache before first refresh
- **WHEN** `CachedCollector.Update` is called before `Start` has been called
- **THEN** it returns no metrics and no error

#### Scenario: Error propagation from cache
- **WHEN** the inner collector's last refresh returned an error
- **THEN** `CachedCollector.Update` returns that cached error alongside any partial metrics

#### Scenario: Graceful shutdown via context
- **WHEN** the context passed to `Start` is cancelled
- **THEN** the background refresh loop exits cleanly

### Requirement: Tiered caching
Sub-collectors SHALL be assigned to one of three tiers based on data volatility. The tier determines whether a sub-collector runs live during scrapes or is wrapped in a `CachedCollector`.

#### Scenario: Live tier (sessions, tasks, queues)
- **WHEN** a sub-collector is classified as live
- **THEN** it is NOT wrapped in a `CachedCollector` and runs directly during each scrape

#### Scenario: Warm tier (counts, activity, calendar)
- **WHEN** a sub-collector is classified as warm
- **THEN** it is wrapped in a `CachedCollector` with the warm cache interval (default 10m)

#### Scenario: Cold tier (system info, library sizes, plugins)
- **WHEN** a sub-collector is classified as cold
- **THEN** it is wrapped in a `CachedCollector` with the cold cache interval (default 30m)

### Requirement: StartCaches lifecycle
`ServiceCollector.StartCaches(ctx)` SHALL launch background refresh goroutines for all `CachedCollector` instances among the registered sub-collectors. The goroutines SHALL exit when the context is cancelled.

#### Scenario: Start caches on application startup
- **WHEN** `StartCaches` is called with the application context
- **THEN** all `CachedCollector` instances begin their background refresh loops

### Requirement: Optional sub-collector registration
The service constructor SHALL support adding sub-collectors conditionally via functional options.

#### Scenario: Jellyfin playback collector conditionally included
- **WHEN** the playback reporting plugin is detected at startup
- **THEN** the playback sub-collector is included in the `ServiceCollector`

#### Scenario: Playback plugin not available
- **WHEN** the playback reporting plugin probe fails
- **THEN** the playback sub-collector is omitted and no playback metrics are emitted
