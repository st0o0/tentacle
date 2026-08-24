## 1. Config Extension

- [x] 1.1 Add `ScanInterval` field to config struct with `TENTACLE_SCAN_INTERVAL` env var, default 6h, 0 to disable
- [x] 1.2 Write tests for scan interval parsing (default, custom, disabled, invalid)

## 2. Cached Collector Infrastructure

- [x] 2.1 Implement `CachedCollector` struct with background goroutine, sync.RWMutex cache, and prometheus.Collector interface
- [x] 2.2 Implement scan function signature, cache replacement on success, keep-previous on failure
- [x] 2.3 Implement graceful shutdown via context cancellation
- [x] 2.4 Implement scan_success meta-metric
- [x] 2.5 Implement zero-interval disabling (no goroutine started)
- [x] 2.6 Write tests for CachedCollector (startup scan, cache retention on failure, shutdown, disabled mode)

## 3. Jellyfin Client Extension

- [x] 3.1 Add MediaSource and MediaStream types to Jellyfin client (Type, Codec, Width, Height, Container)
- [x] 3.2 Verify GetItems with Fields=MediaSources returns populated MediaSources/MediaStreams
- [x] 3.3 Write tests for MediaSource/MediaStream deserialization

## 4. Media Analysis Collector

- [x] 4.1 Implement scan function that iterates video libraries (movies/tvshows), paginates items with Fields=MediaSources
- [x] 4.2 Implement video codec aggregation from primary video stream per item
- [x] 4.3 Implement audio codec aggregation from primary audio stream per item
- [x] 4.4 Implement resolution bucketing (4k/1080p/720p/480p/sd) from video stream width
- [x] 4.5 Implement container format aggregation from MediaSource.Container
- [x] 4.6 Implement subtitle codec aggregation (count per codec, item can contribute to multiple)
- [x] 4.7 Emit scan meta-metrics (last_completed_timestamp, duration, items_analyzed)
- [x] 4.8 Write tests for media analysis (codec counting, resolution bucketing, subtitle handling, library filtering)

## 5. Integration

- [x] 5.1 Register CachedCollector with media analysis scan function in main setup, gated by ScanInterval > 0
- [x] 5.2 Pass application context to CachedCollector for graceful shutdown

## 6. Verification

- [x] 6.1 Run full test suite and verify all tests pass
- [x] 6.2 Run golangci-lint and fix any issues
- [x] 6.3 Build binary and verify it starts with default scan interval
- [x] 6.4 Build binary and verify it starts with TENTACLE_SCAN_INTERVAL=0 (disabled)
