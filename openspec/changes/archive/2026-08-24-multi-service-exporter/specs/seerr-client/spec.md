## Purpose

HTTP client for the Seerr (Jellyseerr/Overseerr) REST API, handling API key authentication, request construction, and response decoding.

## ADDED Requirements

### Requirement: Authentication
All requests SHALL include the `X-Api-Key` header with the configured API key.

#### Scenario: API key in header
- **WHEN** any API request is made
- **THEN** the `X-Api-Key` header contains the configured token

### Requirement: Base URL handling
The client SHALL trim trailing slashes from the base URL.

#### Scenario: Trailing slash
- **WHEN** the base URL is `http://seerr:5055/`
- **THEN** requests target `http://seerr:5055/...` (no double slash)

### Requirement: Context propagation
All API methods SHALL accept a `context.Context` for timeout and cancellation support.

#### Scenario: Timeout
- **WHEN** the context deadline is exceeded
- **THEN** the request is cancelled and an error is returned

### Requirement: Error handling
Non-200 responses SHALL return an error including the status code and up to 512 bytes of the response body.

#### Scenario: 401 response
- **WHEN** the server returns HTTP 401
- **THEN** the error message includes "unexpected status 401"

### Requirement: Status endpoint
`GetStatus()` SHALL call `GET /api/v1/status` and return `Status` (Version, CommitTag).

#### Scenario: Server status
- **WHEN** `/api/v1/status` returns valid JSON
- **THEN** Version and CommitTag are populated

### Requirement: Request Count endpoint
`GetRequestCount()` SHALL call `GET /api/v1/request/count` and return `RequestCount` (Total, Movie, TV, Pending, Approved, Available, Declined).

#### Scenario: Request statistics
- **WHEN** `/api/v1/request/count` returns counts
- **THEN** Total, Pending, Approved, Available, and Declined are populated

### Requirement: Users endpoint
`GetUsers()` SHALL call `GET /api/v1/user` with `take=100&skip=0` and return `UsersResponse` (PageInfo.Pages, PageInfo.Results, Results).

#### Scenario: User list
- **WHEN** `/api/v1/user` returns 5 users
- **THEN** `UsersResponse.PageInfo.Results` is 5
