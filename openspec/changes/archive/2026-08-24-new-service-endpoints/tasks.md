## 1. Arr Client Extensions

- [x] 1.1 Add Backup, Update, BlocklistResponse, DownloadClient types to arr-client
- [x] 1.2 Add GetBackups(), GetUpdates(), GetBlocklist(), GetDownloadClients() methods to arr-client
- [x] 1.3 Modify GetQueue() to use pageSize=250 and populate Records with Status/TrackedDownloadStatus/TrackedDownloadState
- [x] 1.4 Add Application and IndexerStatus types and GetApplications(), GetIndexerStatuses() Prowlarr-specific methods
- [x] 1.5 Write tests for all new arr-client methods

## 2. Seerr Client Extension

- [x] 2.1 Add IssueCount type and GetIssueCount() method to seerr-client
- [x] 2.2 Write tests for GetIssueCount

## 3. Sonarr Collector Extensions

- [x] 3.1 Add `sonarr_queue_by_state_total{state}` metric to queue collector using TrackedDownloadState from queue records
- [x] 3.2 Add `sonarr_backup_total` and `sonarr_backup_latest_timestamp_seconds` metrics
- [x] 3.3 Add `sonarr_update_available{version}` metric
- [x] 3.4 Add `sonarr_blocklist_total` metric
- [x] 3.5 Add `sonarr_download_client_info{name, protocol, priority}` metric
- [x] 3.6 Write tests for all new Sonarr collector metrics

## 4. Radarr Collector Extensions

- [x] 4.1 Add `radarr_queue_by_state_total{state}` metric to queue collector
- [x] 4.2 Add `radarr_backup_total` and `radarr_backup_latest_timestamp_seconds` metrics
- [x] 4.3 Add `radarr_update_available{version}` metric
- [x] 4.4 Add `radarr_blocklist_total` metric
- [x] 4.5 Add `radarr_download_client_info{name, protocol, priority}` metric
- [x] 4.6 Write tests for all new Radarr collector metrics

## 5. Prowlarr Collector Extensions

- [x] 5.1 Add `prowlarr_app_info{name, sync_level, implementation}` and `prowlarr_apps_total` metrics
- [x] 5.2 Add `prowlarr_indexer_disabled{indexer}` metric, resolving indexer ID to name via existing GetIndexers data
- [x] 5.3 Write tests for all new Prowlarr collector metrics

## 6. Seerr Collector Extension

- [x] 6.1 Add `seerr_issues_total` and `seerr_issues_by_status_total{status}` metrics
- [x] 6.2 Write tests for new Seerr collector metrics

## 7. Verification

- [x] 7.1 Run full test suite and verify all tests pass
- [x] 7.2 Run golangci-lint and fix any issues
- [x] 7.3 Build binary and verify it starts without errors
