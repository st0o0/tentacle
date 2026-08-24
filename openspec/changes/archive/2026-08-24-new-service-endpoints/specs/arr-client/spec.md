## MODIFIED Requirements

### Requirement: Queue endpoint
`GetQueue()` SHALL call `GET /api/v3/queue` with `page=1&pageSize=250` query parameters and return `QueueResponse` containing `TotalRecords` and `Records[]` with each record's `Status`, `TrackedDownloadStatus`, and `TrackedDownloadState` fields populated.

#### Scenario: Queue with items and status details
- **WHEN** `/api/v3/queue` returns 5 records
- **THEN** `QueueResponse.TotalRecords` is 5 and each record has Status, TrackedDownloadStatus, and TrackedDownloadState populated

#### Scenario: Empty queue
- **WHEN** `/api/v3/queue` returns 0 records
- **THEN** `QueueResponse.TotalRecords` is 0 and `Records` is an empty slice

## ADDED Requirements

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
