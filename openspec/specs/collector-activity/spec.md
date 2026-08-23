## Purpose

Prometheus collector exposing Jellyfin activity log entry counts by severity and the timestamp of the most recent entry.

## Requirements

### Requirement: Activity log entries by severity
`jellyfin_activity_log_entries_total{severity}` SHALL be a gauge counting the last 100 activity log entries by severity level. Severity labels: Information, Warning, Error.

#### Scenario: Mixed severities
- **WHEN** the last 100 log entries contain 80 Information, 15 Warning, and 5 Error entries
- **THEN** `jellyfin_activity_log_entries_total{severity="Information"}` is 80, `{severity="Warning"}` is 15, `{severity="Error"}` is 5

#### Scenario: No entries of a severity
- **WHEN** there are no Error entries in the last 100
- **THEN** `jellyfin_activity_log_entries_total{severity="Error"}` is 0

### Requirement: Latest entry timestamp
`jellyfin_activity_log_latest_timestamp_seconds` SHALL be a gauge with the Unix timestamp of the most recent activity log entry.

#### Scenario: Recent entry
- **WHEN** the most recent entry is dated 2024-01-15T10:30:00Z
- **THEN** `jellyfin_activity_log_latest_timestamp_seconds` is the Unix timestamp of that date

#### Scenario: No entries
- **WHEN** the activity log is empty
- **THEN** no `jellyfin_activity_log_latest_timestamp_seconds` metric is emitted
