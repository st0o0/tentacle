## Purpose

Jellyfin media analysis collector that scans video/audio/subtitle codec distributions, resolution distributions, and container format distributions across Jellyfin libraries via background scanning.

## Requirements

### Requirement: Video codec distribution
`jellyfin_library_video_codec_total{library, codec}` SHALL be a gauge counting items by primary video stream codec per library. Codec values are lowercase (e.g., "hevc", "h264", "mpeg4", "vp9", "av1").

#### Scenario: Mixed video codecs
- **WHEN** library "Movies" has 250 HEVC and 300 H264 items
- **THEN** `jellyfin_library_video_codec_total{library="Movies", codec="hevc"}` is 250 and `{codec="h264"}` is 300

#### Scenario: Item without video stream
- **WHEN** an item has no video MediaStream
- **THEN** it is not counted in video codec metrics

### Requirement: Audio codec distribution
`jellyfin_library_audio_codec_total{library, codec}` SHALL be a gauge counting items by primary audio stream codec per library. Codec values are lowercase (e.g., "aac", "eac3", "truehd", "flac", "ac3", "dts").

#### Scenario: Mixed audio codecs
- **WHEN** library "Movies" has 280 EAC3 and 200 AAC items
- **THEN** `jellyfin_library_audio_codec_total{library="Movies", codec="eac3"}` is 280 and `{codec="aac"}` is 200

#### Scenario: Multiple audio streams
- **WHEN** an item has multiple audio streams
- **THEN** only the first (primary) audio stream's codec is counted

### Requirement: Resolution distribution
`jellyfin_library_resolution_total{library, resolution}` SHALL be a gauge counting items by resolution bucket per library. Resolution buckets: "4k" (width >= 3840), "1080p" (width >= 1920), "720p" (width >= 1280), "480p" (width >= 720), "sd" (width < 720).

#### Scenario: Mixed resolutions
- **WHEN** library "Movies" has 120 4K, 380 1080p, and 50 720p items
- **THEN** `jellyfin_library_resolution_total{library="Movies", resolution="4k"}` is 120, `{resolution="1080p"}` is 380, `{resolution="720p"}` is 50

#### Scenario: Item without video dimensions
- **WHEN** a video stream has Width=0
- **THEN** it is bucketed as "sd"

### Requirement: Container format distribution
`jellyfin_library_container_total{library, container}` SHALL be a gauge counting items by container format per library. Container values are lowercase (e.g., "mkv", "mp4", "avi", "ts").

#### Scenario: Mixed containers
- **WHEN** library "Movies" has 450 MKV and 100 MP4 items
- **THEN** `jellyfin_library_container_total{library="Movies", container="mkv"}` is 450 and `{container="mp4"}` is 100

#### Scenario: Container from MediaSource
- **WHEN** an item has a MediaSource with Container "mkv"
- **THEN** it is counted under container="mkv"

### Requirement: Subtitle codec distribution
`jellyfin_library_subtitle_codec_total{library, codec}` SHALL be a gauge counting items that have at least one subtitle stream, grouped by subtitle codec per library. Codec values are lowercase (e.g., "srt", "ass", "pgs", "vobsub").

#### Scenario: Mixed subtitle codecs
- **WHEN** library "Movies" has 350 items with SRT subs and 100 with PGS subs
- **THEN** `jellyfin_library_subtitle_codec_total{library="Movies", codec="srt"}` is 350 and `{codec="pgs"}` is 100

#### Scenario: Item with multiple subtitle types
- **WHEN** an item has both SRT and PGS subtitle streams
- **THEN** it is counted once under each codec

#### Scenario: Item without subtitles
- **WHEN** an item has no subtitle streams
- **THEN** it does not contribute to any subtitle codec metric

### Requirement: Scan completion timestamp
`jellyfin_scan_last_completed_timestamp_seconds` SHALL be a gauge with the Unix timestamp of when the last media analysis scan completed.

#### Scenario: After scan
- **WHEN** a scan completes at 2024-06-15T14:00:00Z
- **THEN** `jellyfin_scan_last_completed_timestamp_seconds` is the Unix timestamp of that time

### Requirement: Scan duration
`jellyfin_scan_duration_seconds` SHALL be a gauge with the duration of the last scan in seconds.

#### Scenario: Scan timing
- **WHEN** a scan takes 12.5 seconds
- **THEN** `jellyfin_scan_duration_seconds` is 12.5

### Requirement: Scan items count
`jellyfin_scan_items_analyzed_total` SHALL be a gauge with the total number of items analyzed in the last scan.

#### Scenario: Items analyzed
- **WHEN** a scan processes 560 items across all libraries
- **THEN** `jellyfin_scan_items_analyzed_total` is 560

### Requirement: Library scope
The media analysis SHALL only scan video-type libraries (CollectionType "movies" or "tvshows"). Music, book, and other library types SHALL be skipped.

#### Scenario: Music library skipped
- **WHEN** a library has CollectionType "music"
- **THEN** no media analysis metrics are emitted for that library

#### Scenario: Movie library scanned
- **WHEN** a library has CollectionType "movies"
- **THEN** media analysis metrics are emitted for that library
