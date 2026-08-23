## Purpose

Prometheus collectors exposing Prowlarr metrics including system status, health, indexer inventory, and per-indexer query/grab statistics.

## ADDED Requirements

### Requirement: System up metric
`prowlarr_up` SHALL be a gauge: 1 if the Prowlarr API responds to `/api/v3/system/status`, 0 otherwise.

#### Scenario: Prowlarr reachable
- **WHEN** `/api/v3/system/status` returns HTTP 200
- **THEN** `prowlarr_up` is 1

#### Scenario: Prowlarr unreachable
- **WHEN** `/api/v3/system/status` fails
- **THEN** `prowlarr_up` is 0 and no other system metrics are emitted

### Requirement: System info metric
`prowlarr_system_info{version, branch, runtime}` SHALL be a gauge with constant value 1 exposing Prowlarr metadata.

#### Scenario: System info
- **WHEN** Prowlarr returns version "1.12.0", branch "main", runtime "dotnet"
- **THEN** `prowlarr_system_info{version="1.12.0", branch="main", runtime="dotnet"}` is 1

### Requirement: Health issues metric
`prowlarr_health_issues_total{type, source}` SHALL be a gauge counting health check issues grouped by type and source.

#### Scenario: Health warnings
- **WHEN** Prowlarr reports 1 warning from "IndexerLongTermStatusCheck"
- **THEN** `prowlarr_health_issues_total{type="warning", source="IndexerLongTermStatusCheck"}` is 1

### Requirement: Indexers total metric
`prowlarr_indexers_total` SHALL be a gauge with the total number of configured indexers.

#### Scenario: Indexer count
- **WHEN** Prowlarr has 8 indexers configured
- **THEN** `prowlarr_indexers_total` is 8

### Requirement: Indexers enabled metric
`prowlarr_indexers_enabled_total` SHALL be a gauge with the number of enabled indexers.

#### Scenario: Enabled indexers
- **WHEN** 6 of 8 indexers are enabled
- **THEN** `prowlarr_indexers_enabled_total` is 6

### Requirement: Indexer queries metric
`prowlarr_indexer_queries_total{indexer}` SHALL be a gauge with the total number of queries per indexer.

#### Scenario: Query count
- **WHEN** indexer "NZBgeek" has 500 queries
- **THEN** `prowlarr_indexer_queries_total{indexer="NZBgeek"}` is 500

### Requirement: Indexer grabs metric
`prowlarr_indexer_grabs_total{indexer}` SHALL be a gauge with the total number of successful grabs per indexer.

#### Scenario: Grab count
- **WHEN** indexer "NZBgeek" has 50 grabs
- **THEN** `prowlarr_indexer_grabs_total{indexer="NZBgeek"}` is 50

### Requirement: Indexer failed queries metric
`prowlarr_indexer_failed_queries_total{indexer}` SHALL be a gauge with the total number of failed queries per indexer.

#### Scenario: Failed queries
- **WHEN** indexer "NZBgeek" has 10 failed queries
- **THEN** `prowlarr_indexer_failed_queries_total{indexer="NZBgeek"}` is 10

### Requirement: Indexer failed grabs metric
`prowlarr_indexer_failed_grabs_total{indexer}` SHALL be a gauge with the total number of failed grabs per indexer.

#### Scenario: Failed grabs
- **WHEN** indexer "NZBgeek" has 2 failed grabs
- **THEN** `prowlarr_indexer_failed_grabs_total{indexer="NZBgeek"}` is 2

### Requirement: Indexer average response time metric
`prowlarr_indexer_avg_response_seconds{indexer}` SHALL be a gauge with the average response time in seconds per indexer.

#### Scenario: Response time
- **WHEN** indexer "NZBgeek" has average response time of 1500 milliseconds
- **THEN** `prowlarr_indexer_avg_response_seconds{indexer="NZBgeek"}` is 1.5
