# tentacle

[![CI](https://github.com/st0o0/tentacle/actions/workflows/ci.yml/badge.svg)](https://github.com/st0o0/tentacle/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/st0o0/tentacle?sort=semver)](https://github.com/st0o0/tentacle/releases)
[![GHCR](https://img.shields.io/badge/ghcr.io-st0o0%2Ftentacle-2496ED?logo=docker&logoColor=white)](https://github.com/st0o0/tentacle/pkgs/container/tentacle)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE.md)

Prometheus exporter for your self-hosted media stack. One binary, one
`/metrics` endpoint — covers Jellyfin, Sonarr, Radarr, Prowlarr,
Audiobookshelf, and Seerr. Each service is independently enabled via
environment variables. Single static Go binary, ~scratch image.

Named after the creature that reaches into every corner — tentacle wraps
around your media stack and pulls out every metric worth graphing.

```yaml
services:
  tentacle:
    image: ghcr.io/st0o0/tentacle:latest
    restart: unless-stopped
    environment:
      TENTACLE_JELLYFIN_ADDRESS: http://jellyfin:8096
      TENTACLE_JELLYFIN_TOKEN: ${JELLYFIN_API_TOKEN}
      TENTACLE_SONARR_ADDRESS: http://sonarr:8989
      TENTACLE_SONARR_TOKEN: ${SONARR_API_TOKEN}
      TENTACLE_RADARR_ADDRESS: http://radarr:7878
      TENTACLE_RADARR_TOKEN: ${RADARR_API_TOKEN}
      TENTACLE_PROWLARR_ADDRESS: http://prowlarr:9696
      TENTACLE_PROWLARR_TOKEN: ${PROWLARR_API_TOKEN}
      TENTACLE_AUDIOBOOKSHELF_ADDRESS: http://audiobookshelf:13378
      TENTACLE_AUDIOBOOKSHELF_TOKEN: ${AUDIOBOOKSHELF_API_TOKEN}
      TENTACLE_SEERR_ADDRESS: http://seerr:5055
      TENTACLE_SEERR_TOKEN: ${SEERR_API_TOKEN}
    ports:
      - "9594:9594"
```

Only configure the services you run — a service is enabled when both its
`_ADDRESS` and `_TOKEN` are set.

## Quick start

```bash
docker run --rm \
  -e TENTACLE_JELLYFIN_ADDRESS=http://jellyfin:8096 \
  -e TENTACLE_JELLYFIN_TOKEN=your-api-token \
  -e TENTACLE_SONARR_ADDRESS=http://sonarr:8989 \
  -e TENTACLE_SONARR_TOKEN=your-sonarr-key \
  -p 9594:9594 \
  ghcr.io/st0o0/tentacle:latest
```

Metrics at `:9594/metrics`. Health check at `:9594/healthz`.

## How it works

```
                    ┌─────────────────────────────────┐
                    │           tentacle              │
                    │                                 │
 Jellyfin API ─────▶│  jellyfin collectors (10)       │
 Sonarr API ───────▶│  sonarr collectors (6)          │
 Radarr API ───────▶│  radarr collectors (6)          │──▶ :9594/metrics ──▶ Prometheus
 Prowlarr API ─────▶│  prowlarr collectors (3)        │
 Audiobookshelf ───▶│  audiobookshelf collectors (5)  │
 Seerr API ────────▶│  seerr collectors (4)           │
                    └─────────────────────────────────┘
```

Each collector calls its service's REST API on every Prometheus scrape, converts
the response into gauges, and reports scrape duration and success. No state is
kept between scrapes. All services are read-only.

## Getting API tokens

| Service | How to get the token |
|---|---|
| **Jellyfin** | Dashboard → API Keys → create a new key |
| **Sonarr** | Settings → General → API Key |
| **Radarr** | Settings → General → API Key |
| **Prowlarr** | Settings → General → API Key |
| **Audiobookshelf** | Settings → Users → select user → copy API token |
| **Seerr** | Settings → General → API Key |

## Supported services

### Jellyfin

10 collectors covering sessions, libraries, playback, users, devices, plugins,
tasks, activity, counts, and system health. Metrics use the `jellyfin_` prefix.

| Metric | Type | Labels | Description |
|---|---|---|---|
| `jellyfin_up` | gauge | | Whether Jellyfin is reachable |
| `jellyfin_system_info` | gauge | `version`, `os`, `architecture` | Server information |
| `jellyfin_system_pending_restart` | gauge | | Whether a restart is pending |
| `jellyfin_users_total` | gauge | | Total user count |
| `jellyfin_user_last_login_timestamp_seconds` | gauge | `user` | Last login time |
| `jellyfin_user_last_activity_timestamp_seconds` | gauge | `user` | Last activity time |
| `jellyfin_sessions_active_total` | gauge | | Active session count |
| `jellyfin_sessions_streaming_total` | gauge | `user`, `play_method` | Streaming session count |
| `jellyfin_session_transcoding` | gauge | `user`, `media_type`, `hw_acceleration` | Whether transcoding |
| `jellyfin_session_bitrate_bps` | gauge | `user` | Current bitrate |
| `jellyfin_session_progress_ratio` | gauge | `user`, `media_title` | Playback progress (0-1) |
| `jellyfin_session_paused` | gauge | `user`, `media_title` | Whether paused |
| `jellyfin_session_info` | gauge | `user`, `device`, `client`, `client_version` | Session metadata |
| `jellyfin_session_transcode_info` | gauge | `user`, `video_codec`, `audio_codec`, `container`, `resolution`, `transcode_reason` | Transcode details |
| `jellyfin_session_transcode_framerate` | gauge | `user` | Transcode FPS |
| `jellyfin_session_transcode_completion_ratio` | gauge | `user` | Pre-transcode buffer (0-1) |
| `jellyfin_sessions_bandwidth_total_bps` | gauge | | Total bandwidth |
| `jellyfin_sessions_direct_play_total` | gauge | | Direct play count |
| `jellyfin_sessions_transcode_total` | gauge | | Transcode count |
| `jellyfin_library_items_total` | gauge | `type`, `library`, `collection_type` | Item count per library and type |
| `jellyfin_library_size_bytes` | gauge | `library`, `collection_type` | Total library size |
| `jellyfin_library_latest_added_timestamp_seconds` | gauge | `library`, `collection_type` | Most recently added item |
| `jellyfin_scheduled_task_state` | gauge | `task_name`, `category` | Task state (0=Idle, 1=Running, 2=Cancelling) |
| `jellyfin_scheduled_task_progress_ratio` | gauge | `task_name` | Current progress (0-1) |
| `jellyfin_scheduled_task_last_run_duration_seconds` | gauge | `task_name` | Last run duration |
| `jellyfin_scheduled_task_last_run_success` | gauge | `task_name` | Whether last run succeeded |
| `jellyfin_scheduled_task_last_run_timestamp_seconds` | gauge | `task_name` | Last completion time |
| `jellyfin_activity_log_entries_total` | gauge | `severity` | Entry count by severity |
| `jellyfin_activity_log_entries_by_type_total` | gauge | `type` | Entry count by type |
| `jellyfin_activity_log_latest_timestamp_seconds` | gauge | | Most recent entry |
| `jellyfin_plugins_total` | gauge | | Installed plugin count |
| `jellyfin_plugin_info` | gauge | `name`, `version`, `status` | Plugin details |
| `jellyfin_plugin_update_available` | gauge | `name` | Whether update available |
| `jellyfin_devices_total` | gauge | | Registered device count |
| `jellyfin_device_last_activity_timestamp_seconds` | gauge | `device_name`, `app_name`, `app_version`, `user` | Last device activity |
| `jellyfin_items_total` | gauge | `type` | Global item count (Movie, Series, Episode, Audio, ...) |
| `jellyfin_playback_play_count` | gauge | `user` | Play count per user (last 30 days) |
| `jellyfin_playback_watch_time_seconds` | gauge | `user` | Watch time per user (last 30 days) |

> Playback metrics require the [Playback Reporting](https://github.com/jellyfin/jellyfin-plugin-playback-reporting) plugin.
> The collector is automatically enabled when the plugin is detected at startup.

### Sonarr

6 collectors covering system health, series/episodes, download queue, disk usage,
calendar, and extras (backups, updates, blocklist). Metrics use the `sonarr_` prefix.

| Metric | Type | Labels | Description |
|---|---|---|---|
| `sonarr_up` | gauge | | Whether Sonarr is reachable |
| `sonarr_system_info` | gauge | `version`, `branch`, `runtime` | Server information |
| `sonarr_system_start_time_seconds` | gauge | | Server start timestamp |
| `sonarr_health_issues_total` | gauge | `type`, `source` | Health check issues |
| `sonarr_series_total` | gauge | | Total series count |
| `sonarr_series_monitored_total` | gauge | | Monitored series count |
| `sonarr_series_by_status_total` | gauge | `status` | Series by status |
| `sonarr_seasons_total` | gauge | | Total season count |
| `sonarr_episodes_total` | gauge | | Total episode count |
| `sonarr_episodes_downloaded_total` | gauge | | Downloaded episodes |
| `sonarr_episodes_missing_total` | gauge | | Missing monitored episodes |
| `sonarr_series_size_bytes` | gauge | | Total library size |
| `sonarr_queue_total` | gauge | | Download queue items |
| `sonarr_queue_by_state_total` | gauge | `state` | Queue items by state |
| `sonarr_disk_total_bytes` | gauge | `path` | Total disk space |
| `sonarr_disk_free_bytes` | gauge | `path` | Free disk space |
| `sonarr_calendar_upcoming_total` | gauge | | Episodes airing in next 7 days |
| `sonarr_backup_total` | gauge | | Total backup count |
| `sonarr_backup_latest_timestamp_seconds` | gauge | | Most recent backup |
| `sonarr_update_available` | gauge | `version` | Whether an update is available |
| `sonarr_blocklist_total` | gauge | | Blocked releases count |
| `sonarr_download_client_info` | gauge | `name`, `protocol`, `priority` | Download client details |

### Radarr

6 collectors covering system health, movies, download queue, disk usage,
calendar, and extras (backups, updates, blocklist). Metrics use the `radarr_` prefix.

| Metric | Type | Labels | Description |
|---|---|---|---|
| `radarr_up` | gauge | | Whether Radarr is reachable |
| `radarr_system_info` | gauge | `version`, `branch`, `runtime` | Server information |
| `radarr_system_start_time_seconds` | gauge | | Server start timestamp |
| `radarr_health_issues_total` | gauge | `type`, `source` | Health check issues |
| `radarr_movies_total` | gauge | | Total movie count |
| `radarr_movies_monitored_total` | gauge | | Monitored movies |
| `radarr_movies_by_status_total` | gauge | `status` | Movies by status |
| `radarr_movies_downloaded_total` | gauge | | Movies with files |
| `radarr_movies_missing_total` | gauge | | Missing monitored movies |
| `radarr_movies_size_bytes` | gauge | | Total library size |
| `radarr_queue_total` | gauge | | Download queue items |
| `radarr_queue_by_state_total` | gauge | `state` | Queue items by state |
| `radarr_disk_total_bytes` | gauge | `path` | Total disk space |
| `radarr_disk_free_bytes` | gauge | `path` | Free disk space |
| `radarr_calendar_upcoming_total` | gauge | | Movies releasing in next 30 days |
| `radarr_backup_total` | gauge | | Total backup count |
| `radarr_backup_latest_timestamp_seconds` | gauge | | Most recent backup |
| `radarr_update_available` | gauge | `version` | Whether an update is available |
| `radarr_blocklist_total` | gauge | | Blocked releases count |
| `radarr_download_client_info` | gauge | `name`, `protocol`, `priority` | Download client details |

### Prowlarr

3 collectors covering system health, indexer statistics, and connected apps.
Metrics use the `prowlarr_` prefix.

| Metric | Type | Labels | Description |
|---|---|---|---|
| `prowlarr_up` | gauge | | Whether Prowlarr is reachable |
| `prowlarr_system_info` | gauge | `version`, `branch`, `runtime` | Server information |
| `prowlarr_system_start_time_seconds` | gauge | | Server start timestamp |
| `prowlarr_health_issues_total` | gauge | `type`, `source` | Health check issues |
| `prowlarr_indexers_total` | gauge | | Total indexer count |
| `prowlarr_indexers_enabled_total` | gauge | | Enabled indexers |
| `prowlarr_indexers_by_protocol_total` | gauge | `protocol` | Indexers by protocol |
| `prowlarr_indexer_queries_total` | gauge | `indexer` | Queries per indexer |
| `prowlarr_indexer_grabs_total` | gauge | `indexer` | Grabs per indexer |
| `prowlarr_indexer_failed_queries_total` | gauge | `indexer` | Failed queries |
| `prowlarr_indexer_failed_grabs_total` | gauge | `indexer` | Failed grabs |
| `prowlarr_indexer_avg_response_seconds` | gauge | `indexer` | Average response time |
| `prowlarr_apps_total` | gauge | | Connected application count |
| `prowlarr_app_info` | gauge | `name`, `sync_level`, `implementation` | Application details |
| `prowlarr_indexer_disabled` | gauge | `indexer` | Whether temporarily disabled |

### Audiobookshelf

5 collectors covering system health, libraries, users, sessions, and backups.
Metrics use the `audiobookshelf_` prefix.

| Metric | Type | Labels | Description |
|---|---|---|---|
| `audiobookshelf_up` | gauge | | Whether Audiobookshelf is reachable |
| `audiobookshelf_system_info` | gauge | `version` | Server information |
| `audiobookshelf_libraries_total` | gauge | | Total library count |
| `audiobookshelf_library_items_total` | gauge | `library`, `media_type` | Items per library |
| `audiobookshelf_library_size_bytes` | gauge | `library`, `media_type` | Library size |
| `audiobookshelf_library_duration_seconds` | gauge | `library`, `media_type` | Total content duration |
| `audiobookshelf_library_audio_tracks_total` | gauge | `library`, `media_type` | Audio track count |
| `audiobookshelf_library_authors_total` | gauge | `library`, `media_type` | Author count |
| `audiobookshelf_library_genres_total` | gauge | `library`, `media_type` | Genre count |
| `audiobookshelf_library_missing_total` | gauge | `library`, `media_type` | Items with missing files |
| `audiobookshelf_library_invalid_total` | gauge | `library`, `media_type` | Invalid items |
| `audiobookshelf_library_last_update_timestamp_seconds` | gauge | `library`, `media_type` | Last library update |
| `audiobookshelf_users_total` | gauge | | Total user count |
| `audiobookshelf_users_active_total` | gauge | | Active user count |
| `audiobookshelf_users_online` | gauge | | Online user count |
| `audiobookshelf_user_last_seen_timestamp_seconds` | gauge | `user`, `type` | Last seen time |
| `audiobookshelf_user_listening_time_seconds` | gauge | `user` | Total listening time |
| `audiobookshelf_sessions_total` | gauge | | Total listening sessions |
| `audiobookshelf_backups_total` | gauge | | Total backup count |
| `audiobookshelf_backup_latest_timestamp_seconds` | gauge | | Most recent backup time |

### Seerr

4 collectors covering system health, media requests, users, and issues.
Metrics use the `seerr_` prefix. Works with both Jellyseerr and Overseerr.

| Metric | Type | Labels | Description |
|---|---|---|---|
| `seerr_up` | gauge | | Whether Seerr is reachable |
| `seerr_system_info` | gauge | `version`, `commit` | Server information |
| `seerr_requests_total` | gauge | | Total request count |
| `seerr_requests_by_type_total` | gauge | `media_type` | Requests by type (movie, tv) |
| `seerr_requests_by_status_total` | gauge | `status` | Requests by status |
| `seerr_users_total` | gauge | | Total user count |
| `seerr_issues_total` | gauge | | Total reported issues |
| `seerr_issues_by_status_total` | gauge | `status` | Issues by status |

All collectors also emit `<namespace>_scrape_duration_seconds{collector}` and
`<namespace>_scrape_success{collector}`.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `TENTACLE_JELLYFIN_ADDRESS` | | Jellyfin server URL |
| `TENTACLE_JELLYFIN_TOKEN` | | Jellyfin API token |
| `TENTACLE_SONARR_ADDRESS` | | Sonarr server URL |
| `TENTACLE_SONARR_TOKEN` | | Sonarr API key |
| `TENTACLE_RADARR_ADDRESS` | | Radarr server URL |
| `TENTACLE_RADARR_TOKEN` | | Radarr API key |
| `TENTACLE_PROWLARR_ADDRESS` | | Prowlarr server URL |
| `TENTACLE_PROWLARR_TOKEN` | | Prowlarr API key |
| `TENTACLE_AUDIOBOOKSHELF_ADDRESS` | | Audiobookshelf server URL |
| `TENTACLE_AUDIOBOOKSHELF_TOKEN` | | Audiobookshelf API token |
| `TENTACLE_SEERR_ADDRESS` | | Seerr server URL |
| `TENTACLE_SEERR_TOKEN` | | Seerr API key |
| `TENTACLE_LISTEN_ADDRESS` | `:9594` | Metrics listen address |
| `TENTACLE_SCRAPE_TIMEOUT` | `10s` | Per-collector timeout (seconds or Go duration) |
| `TENTACLE_LOG_LEVEL` | `info` | Log level: debug, info, warn, error |
| `TENTACLE_LOG_FORMAT` | `json` | Log format: json, text |

A service is enabled when both its `_ADDRESS` and `_TOKEN` are set. At least one
service must be configured.

## CLI

```
tentacle              # start the exporter
tentacle version      # print the version and exit
tentacle healthcheck  # check if the server is healthy (exit 0/1)
```

The `healthcheck` subcommand reads `TENTACLE_LISTEN_ADDRESS` to determine the
port, so it works correctly with custom listen addresses. An explicit address can
be passed as an argument: `tentacle healthcheck :8080`.

## Prometheus scrape config

```yaml
scrape_configs:
  - job_name: tentacle
    static_configs:
      - targets: ["tentacle:9594"]
```

## Troubleshooting

**`at least one service must be configured`** — set the address and token
environment variables for at least one service.

**Scrape timeouts** — increase `TENTACLE_SCRAPE_TIMEOUT`. Large libraries
(especially the size calculation) can take longer than the default 10s.

**Playback metrics missing** — install the Playback Reporting plugin in Jellyfin.
The playback collector is auto-detected at startup; restart tentacle after
installing the plugin.

## Development

```bash
go test ./...                                    # tests
go vet ./... && golangci-lint run                 # lint
docker run --rm -i hadolint/hadolint < Dockerfile # Dockerfile lint
docker build -t tentacle:dev .                    # build
```

Commits follow [Conventional Commits](https://www.conventionalcommits.org/).
Releases and the GHCR image are cut automatically by release-please.

## Image

Registry: `ghcr.io/st0o0/tentacle`
Tags: `latest`, `MAJOR.MINOR`, and exact `MAJOR.MINOR.PATCH` per release.
Architectures: `linux/amd64`, `linux/arm64`.

## License

MIT — see [`LICENSE.md`](LICENSE.md).
