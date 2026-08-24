## ADDED Requirements

### Requirement: Scan interval configuration
`TENTACLE_SCAN_INTERVAL` (default `6h`) SHALL configure the interval for background media analysis scans. The value SHALL accept Go duration format. A value of `0` SHALL disable scanning entirely.

#### Scenario: Default scan interval
- **WHEN** `TENTACLE_SCAN_INTERVAL` is not set
- **THEN** `config.ScanInterval` is 6 hours

#### Scenario: Custom interval
- **WHEN** `TENTACLE_SCAN_INTERVAL=12h` is set
- **THEN** `config.ScanInterval` is 12 hours

#### Scenario: Disabled scanning
- **WHEN** `TENTACLE_SCAN_INTERVAL=0` is set
- **THEN** `config.ScanInterval` is 0 and no background scanning occurs

#### Scenario: Invalid duration
- **WHEN** `TENTACLE_SCAN_INTERVAL=banana` is set
- **THEN** config loading returns an error
