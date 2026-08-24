## ADDED Requirements

### Requirement: Queue by state metric
`sonarr_queue_by_state_total{state}` SHALL be a gauge counting queue items grouped by `TrackedDownloadState` (e.g., downloading, importPending, importing, failedPending, warning).

#### Scenario: Mixed queue states
- **WHEN** the queue has 2 downloading, 1 importing, and 1 failedPending items
- **THEN** `sonarr_queue_by_state_total{state="downloading"}` is 2, `{state="importing"}` is 1, `{state="failedPending"}` is 1

#### Scenario: Empty queue
- **WHEN** the queue has no items
- **THEN** no `sonarr_queue_by_state_total` metrics are emitted

### Requirement: Backup metrics
`sonarr_backup_total` SHALL be a gauge with the number of backups. `sonarr_backup_latest_timestamp_seconds` SHALL be a gauge with the Unix timestamp of the most recent backup.

#### Scenario: Backups available
- **WHEN** Sonarr has 3 backups, latest at 2024-06-15T08:00:00Z
- **THEN** `sonarr_backup_total` is 3 and `sonarr_backup_latest_timestamp_seconds` is the Unix timestamp

#### Scenario: No backups
- **WHEN** Sonarr has no backups
- **THEN** `sonarr_backup_total` is 0 and `sonarr_backup_latest_timestamp_seconds` is not emitted

### Requirement: Update available metric
`sonarr_update_available{version}` SHALL be a gauge: 1 if an update is available (latest entry has Installed=false), 0 if up to date. The `version` label SHALL contain the latest available version.

#### Scenario: Update available
- **WHEN** the latest update entry has Version "4.1.0" and Installed=false
- **THEN** `sonarr_update_available{version="4.1.0"}` is 1

#### Scenario: Up to date
- **WHEN** the latest entry has Installed=true, Version "4.0.0"
- **THEN** `sonarr_update_available{version="4.0.0"}` is 0

#### Scenario: No update data
- **WHEN** `/api/v3/update` returns an empty array
- **THEN** no `sonarr_update_available` metric is emitted

### Requirement: Blocklist total metric
`sonarr_blocklist_total` SHALL be a gauge with the total number of blocked releases.

#### Scenario: Blocked releases
- **WHEN** the blocklist has 15 entries
- **THEN** `sonarr_blocklist_total` is 15

### Requirement: Download client info metric
`sonarr_download_client_info{name, protocol, priority}` SHALL be a gauge with value 1 for each configured download client where Enable=true.

#### Scenario: Enabled download client
- **WHEN** client "SABnzbd" with protocol "usenet" and priority 1 is enabled
- **THEN** `sonarr_download_client_info{name="SABnzbd", protocol="usenet", priority="1"}` is 1

#### Scenario: Disabled download client
- **WHEN** a client has Enable=false
- **THEN** no info metric is emitted for that client
