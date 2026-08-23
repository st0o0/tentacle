## Purpose

Prometheus collector exposing installed Jellyfin plugin information and update availability.

## Requirements

### Requirement: Plugin count metric
`jellyfin_plugins_total` SHALL be a gauge with the total number of installed plugins.

#### Scenario: Five plugins installed
- **WHEN** the Jellyfin API returns 5 plugins
- **THEN** `jellyfin_plugins_total` is 5

### Requirement: Plugin info metric
`jellyfin_plugin_info{name, version, status}` SHALL be a gauge with value 1 for each installed plugin.

#### Scenario: Plugin details
- **WHEN** plugin "OpenSubtitles" is installed at version "3.0.0" with status "Active"
- **THEN** `jellyfin_plugin_info{name="OpenSubtitles", version="3.0.0", status="Active"}` is 1

### Requirement: Plugin update metric
`jellyfin_plugin_update_available{name}` SHALL be a gauge: 1 if an update is available, 0 otherwise.

#### Scenario: Update available
- **WHEN** plugin "OpenSubtitles" has `HasUpdate: true`
- **THEN** `jellyfin_plugin_update_available{name="OpenSubtitles"}` is 1

#### Scenario: No update
- **WHEN** plugin "LDAP Auth" has `HasUpdate: false`
- **THEN** `jellyfin_plugin_update_available{name="LDAP Auth"}` is 0
