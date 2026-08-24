## 1. Arr Service Collectors (Sonarr, Radarr, Prowlarr)

- [x] 1.1 Add `sonarr_start_time_seconds` gauge to Sonarr system collector, parsing SystemStatus.StartTime
- [x] 1.2 Add `sonarr_series_by_status_total{status}` gauge vector to Sonarr series collector, counting by Series.Status
- [x] 1.3 Add `sonarr_seasons_total` gauge to Sonarr series collector, summing Series.SeasonCount
- [x] 1.4 Add `radarr_start_time_seconds` gauge to Radarr system collector, parsing SystemStatus.StartTime
- [x] 1.5 Add `radarr_movies_by_status_total{status}` gauge vector to Radarr movies collector, counting by Movie.Status
- [x] 1.6 Add `prowlarr_start_time_seconds` gauge to Prowlarr system collector, parsing SystemStatus.StartTime
- [x] 1.7 Add `prowlarr_indexers_by_protocol_total{protocol}` gauge vector to Prowlarr indexers collector, counting by Indexer.Protocol
- [x] 1.8 Write tests for all new Arr metrics (start time parsing, status grouping, season aggregation, protocol grouping)

## 2. Jellyfin Library Collector

- [x] 2.1 Add `collection_type` label to `jellyfin_library_items_total`, `jellyfin_library_size_bytes`, and `jellyfin_library_latest_added_timestamp_seconds` from VirtualFolder.CollectionType
- [x] 2.2 Update library collector tests for the new label

## 3. Jellyfin Counts Collector

- [x] 3.1 Add `jellyfin_items_total{type="Program"}` and `{type="Item"}` from ItemCounts.ProgramCount and ItemCounts.ItemCount
- [x] 3.2 Update counts collector tests to verify all 12 types are emitted

## 4. Jellyfin Devices Collector

- [x] 4.1 Add `app_version` label to `jellyfin_device_last_activity_timestamp_seconds` from DeviceInfo.AppVersion
- [x] 4.2 Update devices collector tests for the new label

## 5. Jellyfin Activity Collector

- [x] 5.1 Add `jellyfin_activity_log_entries_by_type_total{type}` gauge vector counting entries by ActivityLog.Type
- [x] 5.2 Update activity collector tests for the new type breakdown metric

## 6. Audiobookshelf Collector

- [x] 6.1 Add `audiobookshelf_system_info{version}` gauge using ServerVersion from the most recent backup
- [x] 6.2 Add `audiobookshelf_users_active_total` gauge counting users with IsActive=true
- [x] 6.3 Add `audiobookshelf_library_last_update_timestamp_seconds{library, media_type}` gauge from Library.LastUpdate
- [x] 6.4 Add `audiobookshelf_user_listening_time_seconds{user}` gauge calling GetUserListeningStats per user
- [x] 6.5 Write tests for all new Audiobookshelf metrics (system info, active users, library update, listening stats including N+1 failure handling)

## 7. Verification

- [x] 7.1 Run full test suite and verify all tests pass
- [x] 7.2 Run golangci-lint and fix any issues
- [x] 7.3 Build binary and verify it starts without errors
