## Purpose

HTTP client for the Audiobookshelf REST API, handling Bearer token authentication, request construction, and response decoding. Ported from the Conch project.

## Requirements

### Requirement: Authentication
All requests (except Ping) SHALL include the `Authorization` header with format `Bearer <token>`.

#### Scenario: Token in header
- **WHEN** any authenticated API request is made
- **THEN** the `Authorization` header is `Bearer <configured token>`

### Requirement: Ping endpoint
`Ping()` SHALL call `GET /ping` without authentication and return an error if the status is not 200.

#### Scenario: Server reachable
- **WHEN** `/ping` returns HTTP 200
- **THEN** `Ping()` returns nil

#### Scenario: Server unreachable
- **WHEN** `/ping` fails or returns non-200
- **THEN** `Ping()` returns an error

### Requirement: Base URL handling
The client SHALL trim trailing slashes from the base URL.

#### Scenario: Trailing slash
- **WHEN** the base URL is `http://abs:13378/`
- **THEN** requests target `http://abs:13378/...` (no double slash)

### Requirement: Context propagation
All API methods SHALL accept a `context.Context` for timeout and cancellation support.

#### Scenario: Timeout
- **WHEN** the context deadline is exceeded
- **THEN** the request is cancelled and an error is returned

### Requirement: Error handling
Non-200 responses SHALL return an error including the status code and up to 512 bytes of the response body.

#### Scenario: 403 response
- **WHEN** the server returns HTTP 403
- **THEN** the error message includes "unexpected status 403"

### Requirement: Libraries endpoint
`GetLibraries()` SHALL call `GET /api/libraries` and return a slice of `Library` (Id, Name, MediaType).

#### Scenario: Multiple libraries
- **WHEN** `/api/libraries` returns 2 libraries
- **THEN** a slice of 2 `Library` structs is returned

### Requirement: Library Stats endpoint
`GetLibraryStats()` SHALL call `GET /api/libraries/{id}/stats` and return `LibraryStats` (TotalItems, TotalSize, TotalDuration, TotalAudioTracks, NumAuthors, NumGenres, NumMissing, NumInvalid).

#### Scenario: Library statistics
- **WHEN** `GetLibraryStats(ctx, "lib-id")` is called
- **THEN** all stat fields are populated

### Requirement: Users endpoint
`GetUsers()` SHALL call `GET /api/users` and return a slice of `User` (Id, Username, Type, IsActive, LastSeen).

#### Scenario: User list
- **WHEN** `/api/users` returns 3 users
- **THEN** a slice of 3 `User` structs is returned

### Requirement: Online Users endpoint
`GetOnlineUsers()` SHALL call `GET /api/users/online` and return a slice of `OnlineUser`.

#### Scenario: Online users
- **WHEN** 2 users are online
- **THEN** a slice of 2 `OnlineUser` structs is returned

### Requirement: Sessions endpoint
`GetSessions()` SHALL call `GET /api/sessions` and return `SessionsResponse` (Total, Sessions).

#### Scenario: Session count
- **WHEN** `/api/sessions` returns Total=15
- **THEN** `SessionsResponse.Total` is 15

### Requirement: Backups endpoint
`GetBackups()` SHALL call `GET /api/backups` and return `BackupsResponse` containing a slice of `Backup` (Id, CreatedAt).

#### Scenario: Backup list
- **WHEN** `/api/backups` returns 3 backups
- **THEN** a slice of 3 `Backup` structs is returned with CreatedAt as millisecond timestamps
