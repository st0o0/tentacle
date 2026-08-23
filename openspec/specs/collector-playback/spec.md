## Purpose

Prometheus collector exposing per-user playback statistics from the Jellyfin PlaybackReporting plugin (last 30 days).

## Requirements

### Requirement: PlaybackReporting plugin dependency
This collector requires the PlaybackReporting plugin. If the plugin is not installed, the API call SHALL fail gracefully, logging at debug level (not error).

#### Scenario: Plugin not installed
- **WHEN** `GET /user_usage_stats/user_activity` returns an error
- **THEN** the error is logged at debug level and `jellyfin_scrape_success{collector="playback"}` is 0

### Requirement: Play count metric
`jellyfin_playback_play_count{user}` SHALL be a gauge with the total play count per user over the last 30 days, aggregated across all daily entries.

#### Scenario: User with plays
- **WHEN** user "alice" has daily entries with PlayCount 5, 3, and 2
- **THEN** `jellyfin_playback_play_count{user="alice"}` is 10

### Requirement: Watch time metric
`jellyfin_playback_watch_time_seconds{user}` SHALL be a gauge with the total watch time in seconds per user over the last 30 days, aggregated across all daily entries.

#### Scenario: User watch time
- **WHEN** user "alice" has daily WatchTime entries of 3600, 1800, and 900 seconds
- **THEN** `jellyfin_playback_watch_time_seconds{user="alice"}` is 6300

### Requirement: Empty username filtering
Entries with an empty `UserName` SHALL be skipped.

#### Scenario: Anonymous playback
- **WHEN** an activity entry has an empty `UserName`
- **THEN** it is not included in any playback metrics
