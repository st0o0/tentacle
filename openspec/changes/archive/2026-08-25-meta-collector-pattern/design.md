## Context

Each service (jellyfin, sonarr, radarr, etc.) has multiple sub-collectors that individually implement `prometheus.Collector`. All sub-collectors within a service share `*prometheus.Desc` pointers for scrape instrumentation via a package-level `var scrape = collector.NewScrapeDescs(namespace)`. When multiple collectors send the same descriptor through `Describe()`, the registry panics on duplicate detection.

The node_exporter solves this with a meta-collector pattern: one `prometheus.Collector` per group of sub-collectors, where sub-collectors implement a simpler internal interface.

## Goals / Non-Goals

**Goals:**
- Fix the duplicate descriptor panic so tentacle starts reliably
- Centralize scrape instrumentation (duration/success) in one place per service
- Reduce boilerplate in sub-collectors (no more `Describe()`, no scrape metric emission)
- Keep the same metric output — zero observable change for Prometheus scrapers

**Non-Goals:**
- Concurrent sub-collector scraping (sequential is fine; can add later)
- Changing metric names, labels, or semantics
- Refactoring the test harness beyond what's necessary for the new interface

## Decisions

### 1. SubCollector interface in `internal/collector/collector.go`

```go
type SubCollector interface {
    Name() string
    Describe(ch chan<- *prometheus.Desc)
    Update(ch chan<- prometheus.Metric) error
}
```

Sub-collectors keep their own `Describe()` for domain-specific descriptors (e.g., `jellyfin_system_info`). The wrapper calls each sub-collector's `Describe()` but only the wrapper describes the shared scrape descriptors. `Update()` replaces `Collect()` — it returns an error instead of handling scrape metrics itself.

**Why keep Describe on SubCollector:** Preserves static descriptor validation for domain metrics. The wrapper aggregates all descriptors, so the registry sees them from a single `prometheus.Collector`. This is closer to node_exporter's approach than making everything unchecked.

**Alternative considered:** Empty `Describe()` (unchecked collector). Simpler but loses all registration-time validation. Not worth the trade-off when we can keep it cheaply.

### 2. Generic `ServiceCollector[C]` wrapper in `internal/collector/collector.go`

```go
type ServiceCollector struct {
    scrape       ScrapeDescs
    logger       *slog.Logger
    collectors   []SubCollector
}

func NewServiceCollector(namespace string, logger *slog.Logger, subs ...SubCollector) *ServiceCollector
```

A single concrete struct, not one wrapper per service. Each service package provides a constructor that builds the `ServiceCollector` with its sub-collectors:

```go
// internal/collector/jellyfin/collector.go
func NewCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger, opts ...Option) *collector.ServiceCollector
```

**Why a shared wrapper type:** All six services have identical orchestration logic — iterate subs, time them, emit scrape metrics. Duplicating a wrapper per service would be pure boilerplate.

**Alternative considered:** Per-service wrapper types (like node_exporter's `NodeCollector`). More flexible but tentacle doesn't need per-service custom orchestration logic.

### 3. Optional sub-collectors via functional options

The jellyfin playback collector requires a runtime probe to determine availability. This is handled via an `Option` function:

```go
type Option func(*collectorOptions)

func WithPlayback(enabled bool) Option
```

The probe stays in `main.go` but the result is passed as an option to the jellyfin constructor. This keeps the service package unaware of startup orchestration.

### 4. Sub-collector conversion pattern

Each existing collector changes from:

```go
// Before
func (c *CountsCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.itemsTotal
    ch <- scrape.Duration   // ← removed
    ch <- scrape.Success    // ← removed
}

func (c *CountsCollector) Collect(ch chan<- prometheus.Metric) {
    start := time.Now()
    // ...
    ch <- prometheus.MustNewConstMetric(scrape.Duration, ...)  // ← removed
    ch <- prometheus.MustNewConstMetric(scrape.Success, ...)   // ← removed
}
```

To:

```go
// After
func (c *countsCollector) Name() string { return "counts" }

func (c *countsCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.itemsTotal
}

func (c *countsCollector) Update(ch chan<- prometheus.Metric) error {
    // ... emit domain metrics only, return error on failure
}
```

Sub-collector types become unexported (lowercase) since they're only instantiated by the service constructor.

### 5. `main.go` registration simplification

```go
// Before: 10 individual registrations
reg.MustRegister(
    jellyfincollector.NewSystemCollector(client, timeout, logger),
    jellyfincollector.NewUsersCollector(client, timeout, logger),
    // ... 8 more
)

// After: 1 registration
reg.MustRegister(jellyfincollector.NewCollector(client, timeout, logger, opts...))
```

## Risks / Trade-offs

**[Slight loss of timing granularity]** → Scrape duration now measures only `Update()`, not `Describe()`. In practice `Describe()` is trivial (channel sends), so this is negligible.

**[Sub-collectors become unexported]** → Tests that directly instantiate sub-collectors need updating. Mitigated by testing through the `ServiceCollector` or by keeping test helpers that wrap `Update()`.

**[Sequential execution stays]** → Sub-collectors run sequentially within a service. If one sub-collector is slow, it blocks the rest. This matches current behavior and is acceptable for the scrape volumes tentacle handles. Can add concurrency later if needed.
