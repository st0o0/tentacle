## Why

The first change (harvest-existing-fields) maximizes metrics from already-fetched data. This change adds new API endpoints to the clients, unlocking metrics that require data not currently fetched: queue status breakdowns, backup monitoring, update availability, blocklist tracking, download client status, Prowlarr app sync, and Seerr issue tracking. These are the highest-value monitoring gaps remaining.

## What Changes

- **Sonarr/Radarr**: Fetch full queue records for status/state breakdown; add backup, update, blocklist, and download client endpoints
- **Prowlarr**: Add application sync status and indexer failure status endpoints
- **Seerr**: Add issue count endpoint for issue tracking metrics
- **Arr client**: New shared methods for backup, update, blocklist, download client; modify GetQueue to fetch full records

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `arr-client`: Add GetBackups, GetUpdates, GetBlocklist, GetDownloadClients methods; modify GetQueue to return full records with status fields
- `collector-sonarr`: Add queue_by_status, backup, update_available, blocklist, download_client metrics
- `collector-radarr`: Add queue_by_status, backup, update_available, blocklist, download_client metrics
- `collector-prowlarr`: Add connected_apps and indexer_status metrics
- `collector-seerr`: Add issue count and issues_by_status metrics
- `seerr-client`: Add GetIssueCount method

## Impact

- **Clients**: arr-client gains 5 new methods + 1 modified; seerr-client gains 1 new method
- **Collectors**: 4 collectors extended with new metrics
- **API load**: ~5 additional API calls per Arr service per scrape; 1 additional call for Seerr
- **Config**: No changes

## Non-Goals

- No new collector files
- No config changes
- No Audiobookshelf or Jellyfin changes (those are covered by harvest-existing-fields)
