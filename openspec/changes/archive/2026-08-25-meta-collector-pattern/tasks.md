## 1. Core infrastructure

- [x] 1.1 Add `SubCollector` interface and `ServiceCollector` wrapper to `internal/collector/collector.go`
- [x] 1.2 Add `ServiceCollector` tests verifying Describe deduplication, Collect orchestration, and error isolation

## 2. Convert jellyfin collectors

- [x] 2.1 Add `NewCollector` constructor with `WithPlayback` option to `internal/collector/jellyfin/collector.go`
- [x] 2.2 Convert all jellyfin sub-collectors to `SubCollector` interface (unexport types, remove Describe scrape descs, rename Collect to Update returning error)
- [x] 2.3 Update jellyfin collector tests to use `ServiceCollector` or test `Update` directly

## 3. Convert arr-based collectors

- [x] 3.1 Add `NewCollector` constructor to sonarr, radarr, and prowlarr packages
- [x] 3.2 Convert all sonarr/radarr/prowlarr sub-collectors to `SubCollector` interface
- [x] 3.3 Update sonarr/radarr/prowlarr collector tests

## 4. Convert remaining collectors

- [x] 4.1 Add `NewCollector` constructor to audiobookshelf and seerr packages
- [x] 4.2 Convert all audiobookshelf/seerr sub-collectors to `SubCollector` interface
- [x] 4.3 Update audiobookshelf/seerr collector tests

## 5. Simplify main.go and verify

- [x] 5.1 Update `cmd/tentacle/main.go` to register one `ServiceCollector` per service, move playback probe into `WithPlayback` option
- [x] 5.2 Remove the package-level `var scrape` from each service's `collector.go`
- [x] 5.3 Run full test suite and linter, verify no duplicate descriptor panic
