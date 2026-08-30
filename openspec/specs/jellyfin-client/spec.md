## Purpose

HTTP client for the Jellyfin REST API, handling authentication, request construction, and response decoding for all API endpoints used by the collectors.
## Requirements
### Requirement: Package location
The Jellyfin client SHALL be located at `internal/client/jellyfin/` instead of `internal/jellyfin/`. The public API (types, methods, constructor) SHALL remain unchanged.

#### Scenario: Import path change
- **WHEN** collectors import the Jellyfin client
- **THEN** they import from `github.com/st0o0/tentacle/internal/client/jellyfin`

### Requirement: Authentication
All requests SHALL include the `Authorization` header with format `MediaBrowser Token="<token>"`.

#### Scenario: Token in header
- **WHEN** any API request is made
- **THEN** the `Authorization` header is `MediaBrowser Token="<configured token>"`

### Requirement: Base URL handling
The client SHALL trim trailing slashes from the base URL to prevent double-slash paths.

#### Scenario: Trailing slash
- **WHEN** the base URL is `http://jellyfin:8096/`
- **THEN** requests target `http://jellyfin:8096/System/Info` (no double slash)

### Requirement: Context propagation
All API methods SHALL accept a `context.Context` for timeout and cancellation support.

#### Scenario: Timeout
- **WHEN** the context deadline is exceeded during a request
- **THEN** the request is cancelled and an error is returned

### Requirement: Error handling
Non-200 responses SHALL return an error including the status code and up to 512 bytes of the response body.

#### Scenario: 401 response
- **WHEN** the Jellyfin server returns HTTP 401
- **THEN** the error message includes "unexpected status 401" and the response body excerpt

### Requirement: System Info endpoint
`GetSystemInfo()` SHALL call `GET /System/Info` and decode the response into `SystemInfo` (Version, OperatingSystem, SystemArchitecture, HasPendingRestart).

#### Scenario: Successful system info
- **WHEN** `/System/Info` returns valid JSON
- **THEN** all fields (Version, OperatingSystem, SystemArchitecture, HasPendingRestart) are populated

### Requirement: Users endpoint
`GetUsers()` SHALL call `GET /Users` and return a slice of `User` (Id, Name, LastLoginDate, LastActivityDate).

#### Scenario: Multiple users
- **WHEN** `/Users` returns 3 users
- **THEN** a slice of 3 `User` structs is returned

### Requirement: Sessions endpoint
`GetSessions()` SHALL call `GET /Sessions` and return a slice of `Session` including nested `NowPlayingItem`, `PlayState`, and `TranscodingInfo`.

#### Scenario: Session with playback
- **WHEN** a session has an active `NowPlayingItem` and `TranscodingInfo`
- **THEN** all nested structs are decoded including TranscodeReasons slice

### Requirement: Virtual Folders endpoint
`GetVirtualFolders()` SHALL call `GET /Library/VirtualFolders` and return a slice of `VirtualFolder` (Name, CollectionType, ItemId).

#### Scenario: Library with folders
- **WHEN** `/Library/VirtualFolders` returns 2 folders
- **THEN** each folder's Name and ItemId are populated

### Requirement: Item Counts endpoint
`GetItemCounts()` SHALL call `GET /Items/Counts` and return `ItemCounts` with counts for Movie, Series, Episode, Artist, Album, Song, Book, MusicVideo, Trailer, BoxSet.

#### Scenario: Item counts
- **WHEN** `/Items/Counts` returns valid JSON
- **THEN** all count fields are populated

### Requirement: Items endpoint
`GetItems()` SHALL call `GET /Items` with query parameters: ParentId, Recursive=true, IncludeItemTypes, Fields, Limit, StartIndex. It SHALL return `ItemsResponse` with Items slice and TotalRecordCount. The `MediaSource` struct SHALL include a `Size int64` field mapping `json:"Size"`. When `Fields` includes "MediaSources", each item SHALL include its `MediaSources[]` with `Size`, `Container`, and nested `MediaStreams[]`.

#### Scenario: Paginated items
- **WHEN** `GetItems(ctx, parentID, "Movie", "Size", 500, 0)` is called
- **THEN** the request includes `?ParentId=<id>&Recursive=true&IncludeItemTypes=Movie&Fields=Size&Limit=500`

#### Scenario: Items with MediaSources
- **WHEN** `GetItems(ctx, parentID, "Movie", "MediaSources", 500, 0)` is called
- **THEN** each returned item includes MediaSources with MediaStreams containing Codec, Type, Width, Height

#### Scenario: Items with MediaSources including Size
- **WHEN** `GetItems(ctx, parentID, "Movie", "MediaSources", 500, 0)` is called
- **THEN** each returned item includes MediaSources with Size, Container, and MediaStreams

#### Scenario: MediaSource size populated
- **WHEN** an item has a MediaSource with `Size=4294967296`
- **THEN** `MediaSource.Size` is 4294967296

#### Scenario: MediaStream types
- **WHEN** an item has video, audio, and subtitle streams
- **THEN** MediaStreams includes entries with Type "Video", "Audio", and "Subtitle" respectively

### Requirement: Latest Items endpoint
`GetLatestItems()` SHALL call `GET /Items/Latest` with ParentId and Limit query parameters.

#### Scenario: Latest item
- **WHEN** `GetLatestItems(ctx, parentID, 1)` is called
- **THEN** at most 1 item is returned with DateCreated populated

### Requirement: Scheduled Tasks endpoint
`GetScheduledTasks()` SHALL call `GET /ScheduledTasks` and return a slice of `ScheduledTask` including nested `LastExecutionResult`.

#### Scenario: Running task
- **WHEN** a task has State "Running" and CurrentProgressPercentage 50
- **THEN** both fields are decoded correctly

### Requirement: Activity Log endpoint
`GetActivityLog()` SHALL call `GET /System/ActivityLog/Entries` with a Limit query parameter.

#### Scenario: Activity log with limit
- **WHEN** `GetActivityLog(ctx, 100)` is called
- **THEN** the request includes `?Limit=100`

### Requirement: Plugins endpoint
`GetPlugins()` SHALL call `GET /Plugins` and return a slice of `PluginInfo` (Id, Name, Version, Status, HasUpdate).

#### Scenario: Plugin with update
- **WHEN** a plugin has `HasUpdate: true`
- **THEN** the field is decoded as `true`

### Requirement: Devices endpoint
`GetDevices()` SHALL call `GET /Devices` and return `DevicesResponse` with Items and TotalRecordCount.

#### Scenario: Device info
- **WHEN** `/Devices` returns a device with Name, AppName, LastUserName, DateLastActivity
- **THEN** all fields are populated

### Requirement: Playback Activity endpoint
`GetPlaybackActivity()` SHALL call `GET /user_usage_stats/user_activity` with a `days` query parameter. This endpoint requires the PlaybackReporting plugin.

#### Scenario: Playback stats
- **WHEN** `GetPlaybackActivity(ctx, 30)` is called
- **THEN** the request includes `?days=30` and returns per-user PlayCount and WatchTime

