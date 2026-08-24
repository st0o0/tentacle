## ADDED Requirements

### Requirement: Activity log entries by type
`jellyfin_activity_log_entries_by_type_total{type}` SHALL be a gauge counting the last 100 activity log entries by their Type field (e.g., "AuthenticationSucceeded", "SessionStarted", "UserPolicyUpdated").

#### Scenario: Mixed activity types
- **WHEN** the last 100 log entries contain 40 "SessionStarted", 30 "AuthenticationSucceeded", and 30 "PlaybackStart" entries
- **THEN** `jellyfin_activity_log_entries_by_type_total{type="SessionStarted"}` is 40, `{type="AuthenticationSucceeded"}` is 30, `{type="PlaybackStart"}` is 30

#### Scenario: No entries of a type
- **WHEN** no entries have type "UserDeleted" in the last 100
- **THEN** `jellyfin_activity_log_entries_by_type_total{type="UserDeleted"}` is not emitted

#### Scenario: Empty activity log
- **WHEN** the activity log is empty
- **THEN** no `jellyfin_activity_log_entries_by_type_total` metrics are emitted
