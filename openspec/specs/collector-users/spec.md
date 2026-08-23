## Purpose

Prometheus collector exposing Jellyfin user counts and per-user login/activity timestamps.

## Requirements

### Requirement: User count metric
`jellyfin_users_total` SHALL be a gauge with the total number of Jellyfin users.

#### Scenario: Three users
- **WHEN** the Jellyfin API returns 3 users
- **THEN** `jellyfin_users_total` is 3

### Requirement: Last login timestamp
`jellyfin_user_last_login_timestamp_seconds{user}` SHALL be a gauge with the Unix timestamp of each user's last login.

#### Scenario: User with login date
- **WHEN** user "alice" has a `LastLoginDate`
- **THEN** `jellyfin_user_last_login_timestamp_seconds{user="alice"}` is the Unix timestamp

#### Scenario: User without login date
- **WHEN** user "bob" has no `LastLoginDate`
- **THEN** no `jellyfin_user_last_login_timestamp_seconds` metric is emitted for "bob"

### Requirement: Last activity timestamp
`jellyfin_user_last_activity_timestamp_seconds{user}` SHALL be a gauge with the Unix timestamp of each user's last activity.

#### Scenario: User with activity date
- **WHEN** user "alice" has a `LastActivityDate`
- **THEN** `jellyfin_user_last_activity_timestamp_seconds{user="alice"}` is the Unix timestamp

### Requirement: Jellyfin time parsing
Timestamps SHALL be parsed supporting RFC3339Nano, RFC3339, and bare `2006-01-02T15:04:05` formats.

#### Scenario: Bare timestamp
- **WHEN** a timestamp is `2024-01-15T10:30:00`
- **THEN** it is parsed successfully without timezone (UTC assumed)
