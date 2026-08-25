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

#### Scenario: Collect iterates sub-collectors sequentially
- **WHEN** `ServiceCollector.Collect` is called
- **THEN** it calls each sub-collector's `Update`, measures its duration, and emits `<namespace>_scrape_duration_seconds{collector="<name>"}` and `<namespace>_scrape_success{collector="<name>"}` per sub-collector

#### Scenario: Sub-collector failure does not abort others
- **WHEN** one sub-collector's `Update` returns an error
- **THEN** the error is logged, `scrape_success` is 0 for that collector, and remaining sub-collectors are still executed

### Requirement: Optional sub-collector registration
The service constructor SHALL support adding sub-collectors conditionally via functional options.

#### Scenario: Jellyfin playback collector conditionally included
- **WHEN** the playback reporting plugin is detected at startup
- **THEN** the playback sub-collector is included in the `ServiceCollector`

#### Scenario: Playback plugin not available
- **WHEN** the playback reporting plugin probe fails
- **THEN** the playback sub-collector is omitted and no playback metrics are emitted
