## MODIFIED Requirements

### Requirement: Per-library item count
`jellyfin_library_items_total{type, library, collection_type}` SHALL be a gauge with the item count per type per library. The `collection_type` label SHALL be populated from `VirtualFolder.CollectionType` (e.g., "movies", "tvshows", "music", "books", "mixed"). Item types enumerated: Movie, Series, Episode, MusicAlbum, MusicArtist, Audio, Book.

#### Scenario: Movies in a library
- **WHEN** library "Movies" with CollectionType "movies" has 150 items of type "Movie"
- **THEN** `jellyfin_library_items_total{type="Movie", library="Movies", collection_type="movies"}` is 150

#### Scenario: Item type not present
- **WHEN** library "Music" with CollectionType "music" has 0 items of type "Movie"
- **THEN** `jellyfin_library_items_total{type="Movie", library="Music", collection_type="music"}` is 0

### Requirement: Library size
`jellyfin_library_size_bytes{library, collection_type}` SHALL be a gauge with the total size of all items in a library in bytes. The `collection_type` label SHALL be populated from `VirtualFolder.CollectionType`. Size is calculated by paginating through all items in batches of 500 with the `Size` field.

#### Scenario: Library size calculation
- **WHEN** library "Movies" with CollectionType "movies" has 1200 items with varying sizes
- **THEN** `jellyfin_library_size_bytes{library="Movies", collection_type="movies"}` is the sum of all item sizes, fetched across 3 batches

### Requirement: Latest added timestamp
`jellyfin_library_latest_added_timestamp_seconds{library, collection_type}` SHALL be a gauge with the Unix timestamp of the most recently added item in each library. The `collection_type` label SHALL be populated from `VirtualFolder.CollectionType`.

#### Scenario: Recent addition
- **WHEN** the latest item in library "Movies" (CollectionType "movies") was created at 2024-01-15T10:30:00Z
- **THEN** `jellyfin_library_latest_added_timestamp_seconds{library="Movies", collection_type="movies"}` is the Unix timestamp of that date

#### Scenario: Empty library
- **WHEN** a library has no items
- **THEN** no `jellyfin_library_latest_added_timestamp_seconds` metric is emitted for that library

### Requirement: Partial failure handling
If individual API calls fail (item counts, size, or latest items), the collector SHALL continue processing remaining libraries and types, logging warnings for failures, and set `jellyfin_scrape_success{collector="library"}` to 0.

#### Scenario: One type fails
- **WHEN** the `GetItems` call for type "Episode" fails but others succeed
- **THEN** metrics for other types are still emitted and `jellyfin_scrape_success{collector="library"}` is 0
