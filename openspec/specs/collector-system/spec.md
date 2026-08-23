## Purpose

Prometheus collector exposing Jellyfin system information, pending restart status, and reachability.

## Requirements

### Requirement: System info metric
`jellyfin_system_info{version, os, architecture}` SHALL be a gauge with constant value 1 exposing Jellyfin server metadata.

#### Scenario: System info available
- **WHEN** Prometheus scrapes `/metrics`
- **THEN** `jellyfin_system_info{version="10.9.0", os="Linux", architecture="X64"}` has value 1

### Requirement: Pending restart metric
`jellyfin_system_pending_restart` SHALL be a gauge: 1 if Jellyfin has a pending restart, 0 otherwise.

#### Scenario: No pending restart
- **WHEN** `HasPendingRestart` is false
- **THEN** `jellyfin_system_pending_restart` is 0

#### Scenario: Pending restart
- **WHEN** `HasPendingRestart` is true
- **THEN** `jellyfin_system_pending_restart` is 1

### Requirement: Up metric
`jellyfin_up` SHALL be a gauge: 1 if the Jellyfin API is reachable, 0 if the scrape fails.

#### Scenario: Jellyfin reachable
- **WHEN** `GET /System/Info` succeeds
- **THEN** `jellyfin_up` is 1

#### Scenario: Jellyfin unreachable
- **WHEN** `GET /System/Info` fails
- **THEN** `jellyfin_up` is 0 and no other system metrics are emitted
