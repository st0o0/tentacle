## ADDED Requirements

### Requirement: Start time metric
`prowlarr_start_time_seconds` SHALL be a gauge with the Unix timestamp of when the Prowlarr process started, parsed from `SystemStatus.StartTime`.

#### Scenario: Prowlarr start time
- **WHEN** Prowlarr returns StartTime "2024-06-15T08:00:00Z"
- **THEN** `prowlarr_start_time_seconds` is the Unix timestamp of that date

### Requirement: Indexers by protocol metric
`prowlarr_indexers_by_protocol_total{protocol}` SHALL be a gauge counting indexers grouped by Protocol field (usenet, torrent).

#### Scenario: Mixed protocols
- **WHEN** Prowlarr has 4 usenet and 4 torrent indexers
- **THEN** `prowlarr_indexers_by_protocol_total{protocol="usenet"}` is 4 and `{protocol="torrent"}` is 4

#### Scenario: Single protocol
- **WHEN** all indexers are usenet
- **THEN** only `prowlarr_indexers_by_protocol_total{protocol="usenet"}` is emitted
