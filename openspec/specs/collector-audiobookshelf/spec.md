## Purpose

Prometheus collectors exposing Audiobookshelf server metrics including system health, library statistics, user activity, listening sessions, and backups.

## Requirements

### Requirement: System collector
`audiobookshelf_up` SHALL be a gauge: 1 if the Audiobookshelf `/ping` endpoint responds with HTTP 200, 0 otherwise.

#### Scenario: Server reachable
- **WHEN** `/ping` returns HTTP 200
- **THEN** `audiobookshelf_up` is 1

#### Scenario: Server unreachable
- **WHEN** `/ping` fails
- **THEN** `audiobookshelf_up` is 0

### Requirement: Libraries total metric
`audiobookshelf_libraries_total` SHALL be a gauge with the total number of libraries.

#### Scenario: Three libraries
- **WHEN** the server has 3 libraries
- **THEN** `audiobookshelf_libraries_total` is 3

### Requirement: Library items metric
`audiobookshelf_library_items_total{library, media_type}` SHALL be a gauge with the item count per library.

#### Scenario: Audiobook library
- **WHEN** library "Audiobooks" has 150 items with media_type "book"
- **THEN** `audiobookshelf_library_items_total{library="Audiobooks", media_type="book"}` is 150

### Requirement: Library size metric
`audiobookshelf_library_size_bytes{library, media_type}` SHALL be a gauge with the total size in bytes per library.

#### Scenario: Library disk usage
- **WHEN** library "Audiobooks" uses 50GB
- **THEN** `audiobookshelf_library_size_bytes{library="Audiobooks", media_type="book"}` is 53687091200

### Requirement: Library duration metric
`audiobookshelf_library_duration_seconds{library, media_type}` SHALL be a gauge with the total duration of all items per library.

#### Scenario: Total listening time available
- **WHEN** library "Audiobooks" has 1000 hours of content
- **THEN** `audiobookshelf_library_duration_seconds{library="Audiobooks", media_type="book"}` is 3600000

### Requirement: Library audio tracks metric
`audiobookshelf_library_audio_tracks_total{library, media_type}` SHALL be a gauge with the audio track count per library.

#### Scenario: Audio tracks
- **WHEN** library "Audiobooks" has 5000 audio tracks
- **THEN** `audiobookshelf_library_audio_tracks_total{library="Audiobooks", media_type="book"}` is 5000

### Requirement: Library authors metric
`audiobookshelf_library_authors_total{library, media_type}` SHALL be a gauge with the author count per library.

#### Scenario: Author count
- **WHEN** library "Audiobooks" has 200 authors
- **THEN** `audiobookshelf_library_authors_total{library="Audiobooks", media_type="book"}` is 200

### Requirement: Library genres metric
`audiobookshelf_library_genres_total{library, media_type}` SHALL be a gauge with the genre count per library.

#### Scenario: Genre count
- **WHEN** library "Audiobooks" has 30 genres
- **THEN** `audiobookshelf_library_genres_total{library="Audiobooks", media_type="book"}` is 30

### Requirement: Library missing items metric
`audiobookshelf_library_missing_total{library, media_type}` SHALL be a gauge with the count of items with missing files.

#### Scenario: Missing items
- **WHEN** library "Audiobooks" has 2 items with missing files
- **THEN** `audiobookshelf_library_missing_total{library="Audiobooks", media_type="book"}` is 2

### Requirement: Library invalid items metric
`audiobookshelf_library_invalid_total{library, media_type}` SHALL be a gauge with the count of invalid items.

#### Scenario: Invalid items
- **WHEN** library "Audiobooks" has 1 invalid item
- **THEN** `audiobookshelf_library_invalid_total{library="Audiobooks", media_type="book"}` is 1

### Requirement: Users total metric
`audiobookshelf_users_total` SHALL be a gauge with the total user count.

#### Scenario: Five users
- **WHEN** the server has 5 users
- **THEN** `audiobookshelf_users_total` is 5

### Requirement: Users online metric
`audiobookshelf_users_online` SHALL be a gauge with the currently online user count.

#### Scenario: Two users online
- **WHEN** 2 users are currently online
- **THEN** `audiobookshelf_users_online` is 2

### Requirement: User last seen metric
`audiobookshelf_user_last_seen_timestamp_seconds{user, type}` SHALL be a gauge with the unix timestamp of when the user was last seen. Only emitted for users where LastSeen is non-nil.

#### Scenario: User with last seen
- **WHEN** user "alice" of type "user" was last seen at timestamp 1700000000000 (milliseconds)
- **THEN** `audiobookshelf_user_last_seen_timestamp_seconds{user="alice", type="user"}` is 1700000000

#### Scenario: User without last seen
- **WHEN** user "bob" has no LastSeen value
- **THEN** no `audiobookshelf_user_last_seen_timestamp_seconds` metric is emitted for "bob"

### Requirement: Sessions total metric
`audiobookshelf_sessions_total` SHALL be a gauge with the total listening session count.

#### Scenario: Session count
- **WHEN** the server reports 42 total sessions
- **THEN** `audiobookshelf_sessions_total` is 42

### Requirement: Backups total metric
`audiobookshelf_backups_total` SHALL be a gauge with the total backup count.

#### Scenario: Three backups
- **WHEN** the server has 3 backups
- **THEN** `audiobookshelf_backups_total` is 3

### Requirement: Latest backup timestamp metric
`audiobookshelf_backup_latest_timestamp_seconds` SHALL be a gauge with the unix timestamp of the most recent backup. Only emitted if at least one backup exists.

#### Scenario: Latest backup
- **WHEN** the most recent backup was created at timestamp 1700000000000 (milliseconds)
- **THEN** `audiobookshelf_backup_latest_timestamp_seconds` is 1700000000

#### Scenario: No backups
- **WHEN** no backups exist
- **THEN** `audiobookshelf_backup_latest_timestamp_seconds` is not emitted

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
