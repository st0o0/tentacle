## Context

After harvesting existing response fields, the remaining metric gaps require calling new API endpoints. All target services (Sonarr, Radarr, Prowlarr, Seerr) have stable v3/v1 REST APIs with well-documented endpoints. The arr-client already handles authentication, error handling, and JSON decoding — new methods follow the same pattern.

## Goals / Non-Goals

**Goals:**
- Add highest-value monitoring endpoints across Arr and Seerr services
- Maintain the shared arr-client pattern for Sonarr/Radarr/Prowlarr
- Keep API call count increase reasonable (~5 per Arr service, ~1 for Seerr)

**Non-Goals:**
- No Jellyfin or Audiobookshelf changes
- No caching or background scanning
- No config changes

## Decisions

### Decision 1: Modify GetQueue vs add GetQueueDetails

**Choice**: Modify existing `GetQueue()` to accept a `pageSize` parameter, allowing callers to request full records when status breakdown is needed.

**Why**: A separate method would duplicate the queue URL construction. The existing pageSize=1 optimization for "just count" callers can use pageSize=1 explicitly. The queue collector needs all records for status grouping.

**Alternative**: Separate `GetQueueFull()` — rejected as unnecessary duplication since the same endpoint serves both use cases.

### Decision 2: Queue status grouping granularity

**Choice**: Group by `TrackedDownloadState` (importing, downloading, failedPending, warning, etc.) rather than the higher-level `Status` field.

**Why**: `TrackedDownloadState` gives more actionable categories — "failedPending" is immediately alertable, while the `Status` field is more generic. Both fields are available on each record.

### Decision 3: Backup monitoring scope

**Choice**: Expose backup count and latest backup timestamp only — not individual backup details.

**Why**: Enough for "alert if no backup in X days" without high cardinality. Individual backup info (size, path) adds complexity without clear monitoring value.

### Decision 4: Download client info as info metric

**Choice**: `*_download_client_info{name, protocol, priority}` as a gauge with value 1, rather than detailed per-client metrics.

**Why**: The `/downloadclient` endpoint returns configuration — not runtime status. An info metric lets dashboards show which clients are configured. Runtime download status is already covered by the queue.

### Decision 5: Seerr issues endpoint

**Choice**: Add `GetIssueCount()` calling `/api/v1/issue/count` for aggregate counts rather than paginating `/api/v1/issue`.

**Why**: The count endpoint is a single call returning totals by status. Paginating individual issues would be expensive and high-cardinality.

## Risks / Trade-offs

- **[Increased API load]** → ~5 new calls per Arr service adds latency. Mitigation: All endpoints are lightweight (no pagination needed except queue). Scrape timeout guards against issues.
- **[Queue pagination]** → Large queues (100+ items) mean fetching many records. Mitigation: Use a reasonable pageSize (250) to get all records in one call. Queues rarely exceed this.
- **[Arr API version coupling]** → All endpoints target `/api/v3`. Mitigation: This is the stable API version used by Sonarr v4, Radarr v5, Prowlarr v1.

## Open Questions

_None._
