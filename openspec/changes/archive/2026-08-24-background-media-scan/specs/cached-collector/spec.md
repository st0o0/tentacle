## ADDED Requirements

### Requirement: Cached collector pattern
A `CachedCollector` SHALL run a scan function in a background goroutine on a configurable interval, cache the resulting `[]prometheus.Metric`, and serve cached metrics on `Collect()` calls. It SHALL implement the `prometheus.Collector` interface.

#### Scenario: Normal operation
- **WHEN** the scan interval is 6h and a Prometheus scrape occurs
- **THEN** the `Collect()` method returns cached metrics without triggering a new scan

#### Scenario: First scan on startup
- **WHEN** the cached collector starts
- **THEN** it SHALL trigger the first scan immediately in the background goroutine

#### Scenario: Before first scan completes
- **WHEN** `Collect()` is called before the first scan finishes
- **THEN** no metrics are returned (empty collection)

### Requirement: Cache update on successful scan
The cached collector SHALL replace its cached metrics atomically (under write lock) only on successful scan completion.

#### Scenario: Successful scan
- **WHEN** the scan function returns new metrics without error
- **THEN** the cache is replaced with the new metrics

#### Scenario: Failed scan with existing cache
- **WHEN** the scan function returns an error and a previous cache exists
- **THEN** the previous cached metrics continue to be served

#### Scenario: Failed scan without existing cache
- **WHEN** the scan function returns an error and no previous cache exists
- **THEN** `Collect()` returns no metrics

### Requirement: Scan success metric
The cached collector SHALL emit a `<namespace>_scan_success` gauge: 1 after a successful scan, 0 after a failed scan.

#### Scenario: Successful scan
- **WHEN** the scan completes successfully
- **THEN** `jellyfin_scan_success` is 1

#### Scenario: Failed scan
- **WHEN** the scan fails
- **THEN** `jellyfin_scan_success` is 0

### Requirement: Graceful shutdown
The cached collector SHALL stop its background goroutine when a provided context is cancelled.

#### Scenario: Context cancellation
- **WHEN** the application context is cancelled (SIGTERM)
- **THEN** the background scan goroutine exits cleanly without blocking

### Requirement: Scan interval of zero disables scanning
When the scan interval is 0, the cached collector SHALL NOT start a background goroutine and SHALL return no metrics.

#### Scenario: Disabled scanning
- **WHEN** the scan interval is 0
- **THEN** no background goroutine is started and `Collect()` always returns empty
