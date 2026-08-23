## Purpose

Prometheus collector exposing registered Jellyfin device counts and per-device last activity timestamps.

## Requirements

### Requirement: Device count metric
`jellyfin_devices_total` SHALL be a gauge with the total number of registered devices (from `TotalRecordCount`).

#### Scenario: Ten devices
- **WHEN** the Jellyfin API returns `TotalRecordCount: 10`
- **THEN** `jellyfin_devices_total` is 10

### Requirement: Device last activity metric
`jellyfin_device_last_activity_timestamp_seconds{device_name, app_name, user}` SHALL be a gauge with the Unix timestamp of each device's last activity.

#### Scenario: Device with activity
- **WHEN** device "Fire TV" using app "Jellyfin Android" by user "alice" was last active at 2024-01-15T10:30:00Z
- **THEN** `jellyfin_device_last_activity_timestamp_seconds{device_name="Fire TV", app_name="Jellyfin Android", user="alice"}` is the Unix timestamp

#### Scenario: Device without activity date
- **WHEN** a device has no `DateLastActivity`
- **THEN** no last activity metric is emitted for that device
