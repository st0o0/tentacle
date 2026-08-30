## MODIFIED Requirements

### Requirement: Sonarr Series endpoint
`GetSeries()` SHALL call `GET /api/v3/series` and return a slice of `Series`. The `Series` struct SHALL include a nested `Statistics` struct mapping the `statistics` JSON object, containing `SeasonCount`, `EpisodeCount`, `EpisodeFileCount`, `TotalEpisodeCount`, and `SizeOnDisk`. The top-level fields (SeasonCount, EpisodeCount, EpisodeFileCount, TotalEpisodeCount, SizeOnDisk) SHALL remain mapped for backwards compatibility.

#### Scenario: Multiple series
- **WHEN** `/api/v3/series` returns 50 series
- **THEN** a slice of 50 `Series` structs is returned

#### Scenario: Statistics sub-object populated
- **WHEN** a series has `statistics.seasonCount=5`, `statistics.totalEpisodeCount=52`, `statistics.episodeFileCount=48`, `statistics.sizeOnDisk=123456789`
- **THEN** `Series.Statistics.SeasonCount` is 5, `Series.Statistics.TotalEpisodeCount` is 52, `Series.Statistics.EpisodeFileCount` is 48, `Series.Statistics.SizeOnDisk` is 123456789

#### Scenario: Top-level fields zero with statistics populated
- **WHEN** a series has top-level `sizeOnDisk=0` but `statistics.sizeOnDisk=123456789`
- **THEN** `Series.SizeOnDisk` is 0 and `Series.Statistics.SizeOnDisk` is 123456789
