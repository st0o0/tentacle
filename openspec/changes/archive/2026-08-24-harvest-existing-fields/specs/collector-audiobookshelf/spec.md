## ADDED Requirements

### Requirement: System info metric
`audiobookshelf_system_info{version}` SHALL be a gauge with constant value 1 exposing Audiobookshelf server metadata. The version SHALL be obtained from the `Backup.ServerVersion` field of the most recent backup, or from an alternative source if available. If no version can be determined, the metric SHALL NOT be emitted.

#### Scenario: Version from backup
- **WHEN** the most recent backup has ServerVersion "2.17.0"
- **THEN** `audiobookshelf_system_info{version="2.17.0"}` is 1

#### Scenario: No backups available
- **WHEN** no backups exist
- **THEN** `audiobookshelf_system_info` is not emitted

### Requirement: Active users metric
`audiobookshelf_users_active_total` SHALL be a gauge counting users where `IsActive` is true.

#### Scenario: Mixed active/inactive users
- **WHEN** 4 of 5 users have IsActive=true
- **THEN** `audiobookshelf_users_active_total` is 4

#### Scenario: All users active
- **WHEN** all 3 users have IsActive=true
- **THEN** `audiobookshelf_users_active_total` is 3

### Requirement: Library last update metric
`audiobookshelf_library_last_update_timestamp_seconds{library, media_type}` SHALL be a gauge with the Unix timestamp of when each library was last updated/scanned, from `Library.LastUpdate`.

#### Scenario: Library with update timestamp
- **WHEN** library "Audiobooks" with media_type "book" was last updated at timestamp 1700000000000 (milliseconds)
- **THEN** `audiobookshelf_library_last_update_timestamp_seconds{library="Audiobooks", media_type="book"}` is 1700000000

#### Scenario: Library without update timestamp
- **WHEN** a library has no LastUpdate value
- **THEN** no `audiobookshelf_library_last_update_timestamp_seconds` metric is emitted for that library

### Requirement: User listening time metric
`audiobookshelf_user_listening_time_seconds{user}` SHALL be a gauge with the total listening time in seconds per user, obtained by calling `GetUserListeningStats(userID)` for each user.

#### Scenario: User with listening history
- **WHEN** user "alice" has TotalTime of 86400 seconds
- **THEN** `audiobookshelf_user_listening_time_seconds{user="alice"}` is 86400

#### Scenario: User without listening history
- **WHEN** user "bob" has TotalTime of 0
- **THEN** `audiobookshelf_user_listening_time_seconds{user="bob"}` is 0

#### Scenario: Listening stats API failure
- **WHEN** `GetUserListeningStats` fails for one user
- **THEN** the error is logged as a warning and metrics for other users are still emitted
