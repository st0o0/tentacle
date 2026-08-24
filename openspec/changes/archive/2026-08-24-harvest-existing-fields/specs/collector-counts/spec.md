## MODIFIED Requirements

### Requirement: Global item count metric
`jellyfin_items_total{type}` SHALL be a gauge with the global item count per media type. Types: Movie, Series, Episode, MusicArtist, MusicAlbum, Audio, Book, MusicVideo, Trailer, BoxSet, Program, Item.

#### Scenario: Movie count
- **WHEN** `/Items/Counts` returns `MovieCount: 500`
- **THEN** `jellyfin_items_total{type="Movie"}` is 500

#### Scenario: All types emitted
- **WHEN** `/Items/Counts` returns valid data
- **THEN** all 12 type labels are emitted, including those with count 0

#### Scenario: Program count
- **WHEN** `/Items/Counts` returns `ProgramCount: 25`
- **THEN** `jellyfin_items_total{type="Program"}` is 25

#### Scenario: Item count
- **WHEN** `/Items/Counts` returns `ItemCount: 3000`
- **THEN** `jellyfin_items_total{type="Item"}` is 3000

#### Scenario: API failure
- **WHEN** `GET /Items/Counts` fails
- **THEN** no `jellyfin_items_total` metrics are emitted and `jellyfin_scrape_success{collector="counts"}` is 0
