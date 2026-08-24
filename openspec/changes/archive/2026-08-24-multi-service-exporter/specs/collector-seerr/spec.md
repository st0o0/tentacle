## Purpose

Prometheus collectors exposing Seerr (Jellyseerr/Overseerr) metrics including system status, media request counts by status and type, and user count.

## ADDED Requirements

### Requirement: System up metric
`seerr_up` SHALL be a gauge: 1 if the Seerr API responds to `/api/v1/status`, 0 otherwise.

#### Scenario: Seerr reachable
- **WHEN** `/api/v1/status` returns HTTP 200
- **THEN** `seerr_up` is 1

#### Scenario: Seerr unreachable
- **WHEN** `/api/v1/status` fails
- **THEN** `seerr_up` is 0 and no other metrics are emitted

### Requirement: System info metric
`seerr_system_info{version, commit}` SHALL be a gauge with constant value 1 exposing Seerr metadata.

#### Scenario: System info
- **WHEN** Seerr returns version "2.0.0" and commit "abc1234"
- **THEN** `seerr_system_info{version="2.0.0", commit="abc1234"}` is 1

### Requirement: Requests total metric
`seerr_requests_total` SHALL be a gauge with the total number of media requests.

#### Scenario: Total requests
- **WHEN** Seerr has 100 total requests
- **THEN** `seerr_requests_total` is 100

### Requirement: Requests by media type metric
`seerr_requests_by_type_total{media_type}` SHALL be a gauge with request counts per media type (movie, tv).

#### Scenario: Movie requests
- **WHEN** 60 of 100 requests are for movies
- **THEN** `seerr_requests_by_type_total{media_type="movie"}` is 60

#### Scenario: TV requests
- **WHEN** 40 of 100 requests are for TV
- **THEN** `seerr_requests_by_type_total{media_type="tv"}` is 40

### Requirement: Requests by status metric
`seerr_requests_by_status_total{status}` SHALL be a gauge with request counts per status (pending, approved, available, declined).

#### Scenario: Pending requests
- **WHEN** 5 requests are pending approval
- **THEN** `seerr_requests_by_status_total{status="pending"}` is 5

#### Scenario: Available requests
- **WHEN** 80 requests are available (fulfilled)
- **THEN** `seerr_requests_by_status_total{status="available"}` is 80

### Requirement: Users total metric
`seerr_users_total` SHALL be a gauge with the total number of Seerr users.

#### Scenario: User count
- **WHEN** Seerr has 15 users
- **THEN** `seerr_users_total` is 15
