## MODIFIED Requirements

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
