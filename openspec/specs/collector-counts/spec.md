## Purpose

Prometheus collector exposing global Jellyfin item counts by media type via the `/Items/Counts` API endpoint.

## Requirements

### Requirement: Global item count metric
`jellyfin_items_total{type}` SHALL be a gauge with the global item count per media type. Types: Movie, Series, Episode, MusicArtist, MusicAlbum, Audio, Book, MusicVideo, Trailer, BoxSet.

#### Scenario: Movie count
- **WHEN** `/Items/Counts` returns `MovieCount: 500`
- **THEN** `jellyfin_items_total{type="Movie"}` is 500

#### Scenario: All types emitted
- **WHEN** `/Items/Counts` returns valid data
- **THEN** all 10 type labels are emitted, including those with count 0

#### Scenario: API failure
- **WHEN** `GET /Items/Counts` fails
- **THEN** no `jellyfin_items_total` metrics are emitted and `jellyfin_scrape_success{collector="counts"}` is 0
