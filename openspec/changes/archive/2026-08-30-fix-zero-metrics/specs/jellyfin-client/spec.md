## MODIFIED Requirements

### Requirement: Items endpoint
`GetItems()` SHALL call `GET /Items` with query parameters: ParentId, Recursive=true, IncludeItemTypes, Fields, Limit, StartIndex. It SHALL return `ItemsResponse` with Items slice and TotalRecordCount. The `MediaSource` struct SHALL include a `Size int64` field mapping `json:"Size"`. When `Fields` includes "MediaSources", each item SHALL include its `MediaSources[]` with `Size`, `Container`, and nested `MediaStreams[]`.

#### Scenario: Paginated items
- **WHEN** `GetItems(ctx, parentID, "Movie", "Size", 500, 0)` is called
- **THEN** the request includes `?ParentId=<id>&Recursive=true&IncludeItemTypes=Movie&Fields=Size&Limit=500`

#### Scenario: Items with MediaSources
- **WHEN** `GetItems(ctx, parentID, "Movie", "MediaSources", 500, 0)` is called
- **THEN** each returned item includes MediaSources with MediaStreams containing Codec, Type, Width, Height

#### Scenario: Items with MediaSources including Size
- **WHEN** `GetItems(ctx, parentID, "Movie", "MediaSources", 500, 0)` is called
- **THEN** each returned item includes MediaSources with Size, Container, and MediaStreams

#### Scenario: MediaSource size populated
- **WHEN** an item has a MediaSource with `Size=4294967296`
- **THEN** `MediaSource.Size` is 4294967296

#### Scenario: MediaStream types
- **WHEN** an item has video, audio, and subtitle streams
- **THEN** MediaStreams includes entries with Type "Video", "Audio", and "Subtitle" respectively
