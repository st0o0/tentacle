## 1. Sonarr Statistics Struct

- [x] 1.1 Add `SeriesStatistics` struct to `internal/client/arr/sonarr_types.go` with fields: `SeasonCount`, `EpisodeCount`, `EpisodeFileCount`, `TotalEpisodeCount`, `SizeOnDisk` — mapped to `json:"seasonCount"` etc.
- [x] 1.2 Add `Statistics SeriesStatistics` field to `Series` struct with `json:"statistics"` tag

## 2. Sonarr Collector — Read from Statistics

- [x] 2.1 Update `series.go` loop to read `s.Statistics.TotalEpisodeCount`, `s.Statistics.EpisodeFileCount`, `s.Statistics.SeasonCount`, `s.Statistics.SizeOnDisk` with fallback to top-level when Statistics value is 0

## 3. Jellyfin MediaSource Size Field

- [x] 3.1 Add `Size int64` field with `json:"Size"` to `MediaSource` struct in `internal/client/jellyfin/types.go`

## 4. Jellyfin Library Size — Use MediaSources

- [x] 4.1 Update `librarySize()` in `library.go` to request `Fields=MediaSources` and add `IncludeItemTypes` filter for leaf types (`Movie,Episode,Audio,MusicVideo,Book`)
- [x] 4.2 Change size summation to iterate `item.MediaSources` and sum each `ms.Size` instead of `item.Size`

## 5. Verification

- [x] 5.1 Run `go build ./...` and `golangci-lint run` — no errors
- [x] 5.2 Deploy locally and verify `sonarr_series_size_bytes`, `sonarr_episodes_total`, `sonarr_seasons_total` return non-zero values (manual)
- [x] 5.3 Verify `jellyfin_library_size_bytes` returns non-zero values for all libraries (manual)
