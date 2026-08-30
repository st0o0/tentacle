## MODIFIED Requirements

### Requirement: Episodes total metric
`sonarr_episodes_total` SHALL be a gauge with the total episode count across all series, read from `Series[].Statistics.TotalEpisodeCount`. If `Statistics.TotalEpisodeCount` is 0, the collector SHALL fall back to the top-level `TotalEpisodeCount`.

#### Scenario: Episode count
- **WHEN** all series have a combined 2000 total episodes
- **THEN** `sonarr_episodes_total` is 2000

#### Scenario: Episode count from statistics
- **WHEN** all series have a combined `statistics.totalEpisodeCount` of 2000
- **THEN** `sonarr_episodes_total` is 2000

#### Scenario: Fallback to top-level
- **WHEN** a series has `statistics.totalEpisodeCount=0` and top-level `totalEpisodeCount=50`
- **THEN** the top-level value 50 is used for that series

### Requirement: Episodes downloaded metric
`sonarr_episodes_downloaded_total` SHALL be a gauge with the number of episodes that have files, read from `Series[].Statistics.EpisodeFileCount`. If `Statistics.EpisodeFileCount` is 0, the collector SHALL fall back to the top-level `EpisodeFileCount`.

#### Scenario: Downloaded episodes
- **WHEN** 1800 of 2000 episodes have files
- **THEN** `sonarr_episodes_downloaded_total` is 1800

#### Scenario: Downloaded episodes from statistics
- **WHEN** 1800 of 2000 episodes have files per `statistics.episodeFileCount`
- **THEN** `sonarr_episodes_downloaded_total` is 1800

### Requirement: Seasons total metric
`sonarr_seasons_total` SHALL be a gauge with the total season count aggregated across all series from `Series[].Statistics.SeasonCount`. If `Statistics.SeasonCount` is 0, the collector SHALL fall back to the top-level `SeasonCount`.

#### Scenario: Season count
- **WHEN** Sonarr has 3 series with 5, 3, and 8 seasons respectively
- **THEN** `sonarr_seasons_total` is 16

#### Scenario: Season count from statistics
- **WHEN** Sonarr has 3 series with `statistics.seasonCount` of 5, 3, and 8 respectively
- **THEN** `sonarr_seasons_total` is 16

### Requirement: Series size metric
`sonarr_series_size_bytes` SHALL be a gauge with the total size on disk across all series, read from `Series[].Statistics.SizeOnDisk`. If `Statistics.SizeOnDisk` is 0, the collector SHALL fall back to the top-level `SizeOnDisk`.

#### Scenario: Total library size
- **WHEN** all series consume 1.5TB
- **THEN** `sonarr_series_size_bytes` is 1649267441664

#### Scenario: Total library size from statistics
- **WHEN** all series consume 1.5TB per `statistics.sizeOnDisk`
- **THEN** `sonarr_series_size_bytes` is 1649267441664
