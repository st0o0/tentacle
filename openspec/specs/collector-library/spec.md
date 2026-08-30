## Purpose

Prometheus collector exposing per-library item counts, total size, and latest-added timestamps for Jellyfin virtual folders.
## Requirements
### Requirement: Per-library item count
`jellyfin_library_items_total{type, library, collection_type}` SHALL be a gauge with the item count per type per library. The `collection_type` label SHALL be populated from `VirtualFolder.CollectionType` (e.g., "movies", "tvshows", "music", "books", "mixed"). Item types enumerated: Movie, Series, Episode, MusicAlbum, MusicArtist, Audio, Book.

#### Scenario: Movies in a library
- **WHEN** library "Movies" with CollectionType "movies" has 150 items of type "Movie"
- **THEN** `jellyfin_library_items_total{type="Movie", library="Movies", collection_type="movies"}` is 150

#### Scenario: Item type not present
- **WHEN** library "Music" with CollectionType "music" has 0 items of type "Movie"
- **THEN** `jellyfin_library_items_total{type="Movie", library="Music", collection_type="music"}` is 0

### Requirement: Library size
`jellyfin_library_size_bytes{library, collection_type}` SHALL be a gauge with the total size of all items in a library in bytes. Size SHALL be calculated by paginating through leaf-type items (Movie, Episode, Audio, MusicVideo, Book) with `Fields=MediaSources`, summing `MediaSources[].Size` for each item. The `collection_type` label SHALL be populated from `VirtualFolder.CollectionType`.

#### Scenario: Library size calculation
- **WHEN** library "Movies" with CollectionType "movies" has 1200 items with varying sizes
- **THEN** `jellyfin_library_size_bytes{library="Movies", collection_type="movies"}` is the sum of all `MediaSources[].Size` values, fetched across multiple batches

#### Scenario: Library size from MediaSources
- **WHEN** library "Movies" with CollectionType "movies" has 300 items, each with one MediaSource containing a Size value
- **THEN** `jellyfin_library_size_bytes{library="Movies", collection_type="movies"}` is the sum of all `MediaSources[].Size` values across all items

#### Scenario: Item with multiple MediaSources
- **WHEN** an item has 2 MediaSources with sizes 4GB and 2GB
- **THEN** both sizes (6GB total) are included in the library size sum

#### Scenario: Item with no MediaSources
- **WHEN** an item has an empty MediaSources array
- **THEN** that item contributes 0 to the library size (no error)

#### Scenario: Pagination with leaf types
- **WHEN** library "TV Shows" has 5000 episodes
- **THEN** size is fetched across multiple batches using `IncludeItemTypes=Movie,Episode,Audio,MusicVideo,Book` and each batch sums `MediaSources[].Size`

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

