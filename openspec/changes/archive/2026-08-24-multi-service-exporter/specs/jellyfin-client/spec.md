## MODIFIED Requirements

### Requirement: Package location
The Jellyfin client SHALL be located at `internal/client/jellyfin/` instead of `internal/jellyfin/`. The public API (types, methods, constructor) SHALL remain unchanged.

#### Scenario: Import path change
- **WHEN** collectors import the Jellyfin client
- **THEN** they import from `github.com/st0o0/tentacle/internal/client/jellyfin`
