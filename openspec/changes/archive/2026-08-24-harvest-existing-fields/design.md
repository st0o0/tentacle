## Context

The tentacle exporter has clients for 6 services (Jellyfin, Sonarr, Radarr, Prowlarr, Audiobookshelf, Seerr) with well-defined response types. Analysis revealed ~14 fields already returned by client methods that collectors silently discard. Two Audiobookshelf client methods (`GetStatus()`, `GetUserListeningStats()`) exist but are never called by any collector.

All changes are additive — extending existing collectors with new metrics or labels from data already available.

## Goals / Non-Goals

**Goals:**
- Extract all monitoring-useful data from existing client responses
- Maintain consistency across services (e.g., all Arr services expose start_time)
- Zero increase in API call count for Arr and Jellyfin services

**Non-Goals:**
- No new API endpoints or client response types
- No new collector files
- No config or infrastructure changes
- No changes to scrape interval or caching

## Decisions

### Decision 1: Add labels vs separate metrics for status breakdowns

**Choice**: Use separate `*_by_status_total{status}` gauge vectors rather than adding status labels to existing totals.

**Why**: Adding a `status` label to `sonarr_series_total` would be a breaking change for existing dashboards/alerts. Separate metrics are additive and backwards-compatible.

**Alternative**: Info-style metric `sonarr_series_info{title, status}` per series — rejected due to high cardinality.

### Decision 2: Audiobookshelf GetUserListeningStats N+1 calls

**Choice**: Call `GetUserListeningStats(userID)` for each user returned by `GetUsers()` within the existing users collector.

**Why**: ABS typically has <10 users, so N+1 calls are negligible. The data is only available per-user — there's no bulk endpoint. Adding it to the users collector keeps the user-related data together.

**Risk**: If a deployment has many users (50+), this adds latency. Mitigation: The scrape timeout already guards against this. Log a warning if the call count exceeds a threshold.

### Decision 3: collection_type as label on existing library metrics

**Choice**: Add `collection_type` as an additional label to all `jellyfin_library_*` metrics.

**Why**: This is the most natural dimension for filtering — "show me only movie libraries". The VirtualFolder already carries this field. Since library metrics always have a `library` label, adding `collection_type` doesn't increase cardinality — it's 1:1 with the library name.

**Breaking**: This changes the label set of existing metrics. Prometheus treats different label sets as different time series, so dashboards using these metrics will need to add the new label to queries. However, queries using `{library="Movies"}` without explicit label selection will continue to work.

### Decision 4: Activity log type breakdown scope

**Choice**: Count activity log entries by `Type` field (e.g., "UserAuthenticated", "SessionStarted") from the same 100-entry window already fetched for severity counting.

**Why**: Reuses the same API response — no additional call. The Type field gives more actionable categories than severity alone (e.g., alert on auth failures specifically).

### Decision 5: StartTime parsing for Arr services

**Choice**: Parse `SystemStatus.StartTime` (ISO 8601 string) to Unix timestamp in the system collector for each Arr service.

**Why**: Consistent with how Jellyfin timestamps are handled. Enables uptime calculation in PromQL via `time() - sonarr_start_time_seconds`.

## Risks / Trade-offs

- **[Label change on library metrics]** → Adding `collection_type` label changes time series identity. Existing Grafana dashboards may need query updates. Mitigation: Document in release notes.
- **[Label change on device metrics]** → Adding `app_version` label to device last activity metric. Same mitigation as above.
- **[ABS user listening stats latency]** → N+1 API calls per scrape. Mitigation: Bounded by user count (typically small), guarded by scrape timeout.
- **[Activity log type cardinality]** → Jellyfin has many activity types. Mitigation: The 100-entry window naturally bounds the number of distinct types per scrape.

## Open Questions

_None — all decisions are straightforward extensions of existing patterns._
