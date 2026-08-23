## Purpose

Prometheus collector exposing active Jellyfin session data including streaming counts, transcoding details, bandwidth, and playback progress.

## Requirements

### Requirement: Active sessions count
`jellyfin_sessions_active_total` SHALL be a gauge with the total number of active sessions (including non-streaming).

#### Scenario: Five active sessions
- **WHEN** the Jellyfin API returns 5 sessions
- **THEN** `jellyfin_sessions_active_total` is 5

### Requirement: Streaming count by user and play method
`jellyfin_sessions_streaming_total{user, play_method}` SHALL be a gauge counting sessions currently streaming, grouped by user and play method (DirectPlay, DirectStream, Transcode, Unknown).

#### Scenario: User streaming via DirectPlay
- **WHEN** user "alice" has 2 DirectPlay sessions
- **THEN** `jellyfin_sessions_streaming_total{user="alice", play_method="DirectPlay"}` is 2

#### Scenario: Session without NowPlayingItem
- **WHEN** a session has no `NowPlayingItem`
- **THEN** it is NOT counted in streaming totals

### Requirement: Session info metric
`jellyfin_session_info{user, device, client, client_version}` SHALL be a gauge with value 1 for each streaming session.

#### Scenario: Client info
- **WHEN** user "alice" streams from device "Fire TV" using client "Jellyfin Android" v2.6.0
- **THEN** `jellyfin_session_info{user="alice", device="Fire TV", client="Jellyfin Android", client_version="2.6.0"}` is 1

### Requirement: Transcoding metric
`jellyfin_session_transcoding{user, media_type, hw_acceleration}` SHALL be a gauge: 1 if the video stream is being transcoded, 0 if direct.

#### Scenario: Hardware-accelerated transcoding
- **WHEN** a session is transcoding with hardware acceleration
- **THEN** `jellyfin_session_transcoding{user="alice", media_type="Movie", hw_acceleration="true"}` is 1

#### Scenario: Direct video playback
- **WHEN** `IsVideoDirect` is true
- **THEN** `jellyfin_session_transcoding{..., hw_acceleration="false"}` is 0

### Requirement: Transcode info metric
`jellyfin_session_transcode_info{user, video_codec, audio_codec, container, resolution, transcode_reason}` SHALL be a gauge with value 1 for sessions with transcoding info.

#### Scenario: Transcode details
- **WHEN** a session transcodes to h264/aac in ts at 1920x1080
- **THEN** `jellyfin_session_transcode_info{..., video_codec="h264", audio_codec="aac", container="ts", resolution="1920x1080", transcode_reason="ContainerNotSupported"}` is 1

### Requirement: Bitrate metric
`jellyfin_session_bitrate_bps{user}` SHALL be a gauge with the total bitrate in bits per second, summed across all streaming sessions for a user.

#### Scenario: Multiple streams
- **WHEN** user "alice" has two streams with bitrates 5Mbps and 3Mbps
- **THEN** `jellyfin_session_bitrate_bps{user="alice"}` is 8000000

### Requirement: Transcode framerate metric
`jellyfin_session_transcode_framerate{user}` SHALL be a gauge with the current transcode framerate, summed per user.

#### Scenario: Transcode at 30fps
- **WHEN** a session transcodes at 30fps
- **THEN** `jellyfin_session_transcode_framerate{user="alice"}` is 30

### Requirement: Transcode completion metric
`jellyfin_session_transcode_completion_ratio{user}` SHALL be a gauge from 0 to 1 representing how much of the file has been pre-transcoded.

#### Scenario: Half transcoded
- **WHEN** `CompletionPercentage` is 50
- **THEN** `jellyfin_session_transcode_completion_ratio{user="alice"}` is 0.5

### Requirement: Playback progress metric
`jellyfin_session_progress_ratio{user, media_title}` SHALL be a gauge from 0 to 1 calculated as `PositionTicks / RunTimeTicks`.

#### Scenario: Halfway through movie
- **WHEN** PositionTicks is 50% of RunTimeTicks
- **THEN** `jellyfin_session_progress_ratio{user="alice", media_title="Movie"}` is 0.5

#### Scenario: Zero-length media
- **WHEN** `RunTimeTicks` is 0
- **THEN** no progress metric is emitted for that session

### Requirement: Paused metric
`jellyfin_session_paused{user, media_title}` SHALL be a gauge: 1 if paused, 0 if playing.

#### Scenario: Paused session
- **WHEN** `IsPaused` is true
- **THEN** `jellyfin_session_paused{user="alice", media_title="Movie"}` is 1

### Requirement: Total bandwidth metric
`jellyfin_sessions_bandwidth_total_bps` SHALL be a gauge with the sum of all streaming session bitrates.

#### Scenario: Total bandwidth
- **WHEN** three sessions have bitrates 5M, 3M, and 2M
- **THEN** `jellyfin_sessions_bandwidth_total_bps` is 10000000

### Requirement: Direct play and transcode counts
`jellyfin_sessions_direct_play_total` SHALL count sessions using DirectPlay or DirectStream. `jellyfin_sessions_transcode_total` SHALL count sessions using Transcode.

#### Scenario: Mixed playback methods
- **WHEN** 2 sessions use DirectPlay and 1 uses Transcode
- **THEN** `jellyfin_sessions_direct_play_total` is 2 and `jellyfin_sessions_transcode_total` is 1

### Requirement: Unknown user fallback
Sessions without a `UserName` SHALL use the label value `"unknown"`.

#### Scenario: Anonymous session
- **WHEN** a streaming session has an empty `UserName`
- **THEN** metrics use `user="unknown"`
