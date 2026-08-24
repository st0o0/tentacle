## MODIFIED Requirements

### Requirement: Device last activity metric
`jellyfin_device_last_activity_timestamp_seconds{device_name, app_name, app_version, user}` SHALL be a gauge with the Unix timestamp of each device's last activity. The `app_version` label SHALL be populated from `DeviceInfo.AppVersion`.

#### Scenario: Device with activity and version
- **WHEN** device "Fire TV" using app "Jellyfin Android" version "2.6.0" by user "alice" was last active at 2024-01-15T10:30:00Z
- **THEN** `jellyfin_device_last_activity_timestamp_seconds{device_name="Fire TV", app_name="Jellyfin Android", app_version="2.6.0", user="alice"}` is the Unix timestamp

#### Scenario: Device without version
- **WHEN** a device has no AppVersion
- **THEN** the `app_version` label SHALL be an empty string

#### Scenario: Device without activity date
- **WHEN** a device has no `DateLastActivity`
- **THEN** no last activity metric is emitted for that device
