## Purpose

Shared HTTP client for the Arr API (Sonarr, Radarr, Prowlarr), handling authentication, request construction, and response decoding for `/api/v3` endpoints.
## Requirements
### Requirement: Authentication
All requests SHALL include the `X-Api-Key` header with the configured API key.

#### Scenario: API key in header
- **WHEN** any API request is made
- **THEN** the `X-Api-Key` header contains the configured token

### Requirement: Base URL and API path
The client SHALL prepend `/api/v3` to all endpoint paths. Trailing slashes on the base URL SHALL be trimmed.

#### Scenario: Endpoint path construction
- **WHEN** `GetSystemStatus()` is called with base URL `http://sonarr:8989`
- **THEN** the request targets `http://sonarr:8989/api/v3/system/status`

#### Scenario: Trailing slash in base URL
- **WHEN** the base URL is `http://sonarr:8989/`
- **THEN** requests target `http://sonarr:8989/api/v3/...` (no double slash)

### Requirement: Context propagation
All API methods SHALL accept a `context.Context` for timeout and cancellation support.

#### Scenario: Timeout
- **WHEN** the context deadline is exceeded during a request
- **THEN** the request is cancelled and an error is returned

### Requirement: Error handling
Non-200 responses SHALL return an error including the status code and up to 512 bytes of the response body.

#### Scenario: 401 unauthorized
- **WHEN** the Arr service returns HTTP 401
- **THEN** the error message includes "unexpected status 401" and the response body excerpt

### Requirement: System Status endpoint
`GetSystemStatus()` SHALL call `GET /api/v3/system/status` and return `SystemStatus` (Version, StartTime, Branch, RuntimeName, RuntimeVersion).

#### Scenario: Successful system status
- **WHEN** `/api/v3/system/status` returns valid JSON
- **THEN** Version, StartTime, Branch, RuntimeName, RuntimeVersion are populated

### Requirement: Health endpoint
`GetHealth()` SHALL call `GET /api/v3/health` and return a slice of `HealthCheck` (Source, Type, Message, WikiUrl).

#### Scenario: Health warnings
- **WHEN** `/api/v3/health` returns 2 health checks
- **THEN** a slice of 2 `HealthCheck` structs is returned with Source and Type populated

#### Scenario: Healthy system
- **WHEN** `/api/v3/health` returns an empty array
- **THEN** an empty slice is returned

### Requirement: Queue endpoint
`GetQueue()` SHALL call `GET /api/v3/queue` with `page=1&pageSize=250` query parameters and return `QueueResponse` containing `TotalRecords` and `Records[]` with each record's `Status`, `TrackedDownloadStatus`, and `TrackedDownloadState` fields populated.

#### Scenario: Queue with items and status details
- **WHEN** `/api/v3/queue` returns 5 records
- **THEN** `QueueResponse.TotalRecords` is 5 and each record has Status, TrackedDownloadStatus, and TrackedDownloadState populated

#### Scenario: Empty queue
- **WHEN** `/api/v3/queue` returns 0 records
- **THEN** `QueueResponse.TotalRecords` is 0 and `Records` is an empty slice

### Requirement: Root Folder endpoint
`GetRootFolders()` SHALL call `GET /api/v3/rootfolder` and return a slice of `RootFolder` (Path, FreeSpace, TotalSpace).

#### Scenario: Disk space info
- **WHEN** `/api/v3/rootfolder` returns a folder with FreeSpace and TotalSpace
- **THEN** both values are populated as int64 byte counts

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

### Requirement: Sonarr Wanted Missing endpoint
`GetWantedMissing()` SHALL call `GET /api/v3/wanted/missing` with `page=1&pageSize=1` and return `WantedResponse` (TotalRecords).

#### Scenario: Missing episodes count
- **WHEN** `/api/v3/wanted/missing` returns TotalRecords=12
- **THEN** `WantedResponse.TotalRecords` is 12

### Requirement: Sonarr Calendar endpoint
`GetCalendar()` SHALL call `GET /api/v3/calendar` with `start` and `end` query parameters and return a slice of upcoming episodes or movies.

#### Scenario: Upcoming episodes
- **WHEN** `GetCalendar(ctx, start, end)` is called with a 7-day range
- **THEN** all episodes airing in that range are returned

### Requirement: Radarr Movies endpoint
`GetMovies()` SHALL call `GET /api/v3/movie` and return a slice of `Movie` (Title, Monitored, Status, HasFile, SizeOnDisk, Year).

#### Scenario: Movie library
- **WHEN** `/api/v3/movie` returns 200 movies
- **THEN** a slice of 200 `Movie` structs is returned

### Requirement: Prowlarr Indexers endpoint
`GetIndexers()` SHALL call `GET /api/v3/indexer` and return a slice of `Indexer` (Name, Enable, Protocol, Priority).

#### Scenario: Multiple indexers
- **WHEN** `/api/v3/indexer` returns 5 indexers
- **THEN** a slice of 5 `Indexer` structs is returned with Enable as bool

### Requirement: Prowlarr Indexer Stats endpoint
`GetIndexerStats()` SHALL call `GET /api/v3/indexerstats` and return `IndexerStatsResponse` containing per-indexer stats (NumberOfQueries, NumberOfGrabs, NumberOfFailedQueries, NumberOfFailedGrabs, AverageResponseTime).

#### Scenario: Indexer statistics
- **WHEN** `/api/v3/indexerstats` returns stats for 3 indexers
- **THEN** each indexer's query count, grab count, and average response time are populated

### Requirement: Backups endpoint
`GetBackups()` SHALL call `GET /api/v3/system/backup` and return a slice of `Backup` (Id, Name, Path, Type, Time).

#### Scenario: Backup list
- **WHEN** `/api/v3/system/backup` returns 3 backups
- **THEN** a slice of 3 `Backup` structs is returned with Time as ISO 8601 string

#### Scenario: No backups
- **WHEN** `/api/v3/system/backup` returns an empty array
- **THEN** an empty slice is returned

### Requirement: Updates endpoint
`GetUpdates()` SHALL call `GET /api/v3/update` and return a slice of `Update` (Version, Installed, Latest).

#### Scenario: Update available
- **WHEN** `/api/v3/update` returns an entry with Installed=false and Latest=true
- **THEN** the update's Version indicates an available update

#### Scenario: Already up to date
- **WHEN** the latest entry has Installed=true
- **THEN** no update is pending

### Requirement: Blocklist endpoint
`GetBlocklist()` SHALL call `GET /api/v3/blocklist` with `page=1&pageSize=1` and return `BlocklistResponse` (TotalRecords).

#### Scenario: Blocklist count
- **WHEN** `/api/v3/blocklist` returns TotalRecords=15
- **THEN** `BlocklistResponse.TotalRecords` is 15

### Requirement: Download Clients endpoint
`GetDownloadClients()` SHALL call `GET /api/v3/downloadclient` and return a slice of `DownloadClient` (Name, Protocol, Priority, Enable).

#### Scenario: Multiple download clients
- **WHEN** `/api/v3/downloadclient` returns 2 clients
- **THEN** a slice of 2 `DownloadClient` structs is returned

### Requirement: Prowlarr Applications endpoint
`GetApplications()` SHALL call `GET /api/v3/applications` and return a slice of `Application` (Name, SyncLevel, Implementation).

#### Scenario: Connected applications
- **WHEN** `/api/v3/applications` returns 3 apps
- **THEN** a slice of 3 `Application` structs is returned with SyncLevel (e.g., "fullSync", "disabled")

### Requirement: Prowlarr Indexer Status endpoint
`GetIndexerStatuses()` SHALL call `GET /api/v3/indexerstatus` and return a slice of `IndexerStatus` (IndexerId, DisabledTill).

#### Scenario: Failed indexer
- **WHEN** an indexer has a DisabledTill timestamp in the future
- **THEN** the IndexerStatus shows the indexer is temporarily disabled

#### Scenario: All indexers healthy
- **WHEN** `/api/v3/indexerstatus` returns an empty array
- **THEN** no indexers are disabled

