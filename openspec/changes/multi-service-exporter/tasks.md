## 1. Restructure Package Layout

- [x] 1.1 Move `internal/jellyfin/` to `internal/client/jellyfin/` (client.go, types.go)
- [x] 1.2 Create `internal/collector/jellyfin/` and move all Jellyfin collector files (system.go, users.go, sessions.go, library.go, tasks.go, activity.go, plugins.go, devices.go, counts.go, playback.go) from `internal/collector/`
- [x] 1.3 Extract shared scrape helpers in `internal/collector/collector.go` to accept namespace as parameter (`NewScrapeDescs(namespace)`)
- [x] 1.4 Update all Jellyfin collector imports and namespace references to use new package paths
- [x] 1.5 Update `cmd/tentacle/main.go` imports and collector registration
- [x] 1.6 Verify build compiles and existing behavior is unchanged (`go build ./...`, `go vet ./...`)

## 2. Multi-Service Config

- [x] 2.1 Add `ServiceConfig` struct (Address, Token) and per-service fields to `Config` (Jellyfin, Sonarr, Radarr, Prowlarr, Audiobookshelf, Seerr as `*ServiceConfig`)
- [x] 2.2 Implement per-service env var loading with partial-config validation (address set without token = error)
- [x] 2.3 Add validation that at least one service is configured
- [x] 2.4 Keep global settings (ListenAddress, ScrapeTimeout, LogLevel, LogFormat, DisablePlayback) unchanged
- [x] 2.5 Write config tests for: no services, single service, multiple services, partial config error
- [x] 2.6 Update `main.go` to use conditional service registration based on config

## 3. Shared Arr Client

- [x] 3.1 Create `internal/client/arr/client.go` with `ArrClient` struct (baseURL, apiKey, httpClient), `X-Api-Key` auth, `/api/v3` prefix, shared `get()` method
- [x] 3.2 Create `internal/client/arr/types.go` with shared types: `SystemStatus`, `HealthCheck`, `QueueResponse`, `RootFolder`
- [x] 3.3 Implement shared endpoints: `GetSystemStatus()`, `GetHealth()`, `GetQueue()`, `GetRootFolders()`
- [x] 3.4 Write tests for Arr client (auth header, URL construction, error handling, JSON decoding)

## 4. Sonarr

- [x] 4.1 Add Sonarr-specific methods to Arr client: `GetSeries()`, `GetWantedMissing()`, `GetCalendar()` in `internal/client/arr/sonarr.go` with types in `sonarr_types.go`
- [x] 4.2 Create `internal/collector/sonarr/system.go` — `sonarr_up`, `sonarr_system_info`, `sonarr_health_issues_total`
- [x] 4.3 Create `internal/collector/sonarr/series.go` — `sonarr_series_total`, `sonarr_series_monitored_total`, `sonarr_episodes_total`, `sonarr_episodes_downloaded_total`, `sonarr_episodes_missing_total`, `sonarr_series_size_bytes`
- [x] 4.4 Create `internal/collector/sonarr/queue.go` — `sonarr_queue_total`
- [x] 4.5 Create `internal/collector/sonarr/disk.go` — `sonarr_disk_total_bytes`, `sonarr_disk_free_bytes`
- [x] 4.6 Create `internal/collector/sonarr/calendar.go` — `sonarr_calendar_upcoming_total`
- [x] 4.7 Register Sonarr collectors in `main.go` when config is present
- [x] 4.8 Write tests for Sonarr collectors

## 5. Radarr

- [x] 5.1 Add Radarr-specific methods to Arr client: `GetMovies()`, `GetWantedMissingMovies()`, `GetMovieCalendar()` in `internal/client/arr/radarr.go` with types in `radarr_types.go`
- [x] 5.2 Create `internal/collector/radarr/system.go` — `radarr_up`, `radarr_system_info`, `radarr_health_issues_total`
- [x] 5.3 Create `internal/collector/radarr/movies.go` — `radarr_movies_total`, `radarr_movies_monitored_total`, `radarr_movies_downloaded_total`, `radarr_movies_missing_total`, `radarr_movies_size_bytes`
- [x] 5.4 Create `internal/collector/radarr/queue.go` — `radarr_queue_total`
- [x] 5.5 Create `internal/collector/radarr/disk.go` — `radarr_disk_total_bytes`, `radarr_disk_free_bytes`
- [x] 5.6 Create `internal/collector/radarr/calendar.go` — `radarr_calendar_upcoming_total`
- [x] 5.7 Register Radarr collectors in `main.go` when config is present
- [x] 5.8 Write tests for Radarr collectors

## 6. Prowlarr

- [x] 6.1 Add Prowlarr-specific methods to Arr client: `GetIndexers()`, `GetIndexerStats()` in `internal/client/arr/prowlarr.go` with types in `prowlarr_types.go`
- [x] 6.2 Create `internal/collector/prowlarr/system.go` — `prowlarr_up`, `prowlarr_system_info`, `prowlarr_health_issues_total`
- [x] 6.3 Create `internal/collector/prowlarr/indexers.go` — `prowlarr_indexers_total`, `prowlarr_indexers_enabled_total`, `prowlarr_indexer_queries_total`, `prowlarr_indexer_grabs_total`, `prowlarr_indexer_failed_queries_total`, `prowlarr_indexer_failed_grabs_total`, `prowlarr_indexer_avg_response_seconds`
- [x] 6.4 Register Prowlarr collectors in `main.go` when config is present
- [x] 6.5 Write tests for Prowlarr collectors

## 7. Audiobookshelf

- [x] 7.1 Port `internal/audiobookshelf/` from Conch to `internal/client/audiobookshelf/` (client.go, types.go), adapting to Tentacle conventions
- [x] 7.2 Create `internal/collector/audiobookshelf/system.go` — `audiobookshelf_up`
- [x] 7.3 Create `internal/collector/audiobookshelf/libraries.go` — `audiobookshelf_libraries_total`, `audiobookshelf_library_items_total`, `audiobookshelf_library_size_bytes`, `audiobookshelf_library_duration_seconds`, `audiobookshelf_library_audio_tracks_total`, `audiobookshelf_library_authors_total`, `audiobookshelf_library_genres_total`, `audiobookshelf_library_missing_total`, `audiobookshelf_library_invalid_total`
- [x] 7.4 Create `internal/collector/audiobookshelf/users.go` — `audiobookshelf_users_total`, `audiobookshelf_users_online`, `audiobookshelf_user_last_seen_timestamp_seconds`
- [x] 7.5 Create `internal/collector/audiobookshelf/sessions.go` — `audiobookshelf_sessions_total`
- [x] 7.6 Create `internal/collector/audiobookshelf/backups.go` — `audiobookshelf_backups_total`, `audiobookshelf_backup_latest_timestamp_seconds`
- [x] 7.7 Register Audiobookshelf collectors in `main.go` when config is present
- [x] 7.8 Write tests for Audiobookshelf collectors

## 8. Seerr

- [x] 8.1 Create `internal/client/seerr/client.go` with `X-Api-Key` auth, `/api/v1` prefix, shared `get()` method
- [x] 8.2 Create `internal/client/seerr/types.go` with types: `Status`, `RequestCount`, `UsersResponse`, `PageInfo`
- [x] 8.3 Implement endpoints: `GetStatus()`, `GetRequestCount()`, `GetUsers()`
- [x] 8.4 Create `internal/collector/seerr/system.go` — `seerr_up`, `seerr_system_info`
- [x] 8.5 Create `internal/collector/seerr/requests.go` — `seerr_requests_total`, `seerr_requests_by_type_total`, `seerr_requests_by_status_total`
- [x] 8.6 Create `internal/collector/seerr/users.go` — `seerr_users_total`
- [x] 8.7 Register Seerr collectors in `main.go` when config is present
- [x] 8.8 Write tests for Seerr collectors

## 9. Integration & Finalize

- [x] 9.1 Update `main.go` startup logging to list which services are enabled
- [x] 9.2 Verify full build: `go build ./...`, `go vet ./...`, `golangci-lint run`
- [x] 9.3 Run full test suite: `go test -race ./...`
- [x] 9.4 Update Dockerfile (no functional changes expected, verify build still works)
- [x] 9.5 Update README with new env vars and supported services
