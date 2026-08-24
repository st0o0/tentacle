## ADDED Requirements

### Requirement: Connected applications metric
`prowlarr_app_info{name, sync_level, implementation}` SHALL be a gauge with value 1 for each connected application.

#### Scenario: Connected apps
- **WHEN** Prowlarr has Sonarr (fullSync) and Radarr (fullSync) connected
- **THEN** `prowlarr_app_info{name="Sonarr", sync_level="fullSync", implementation="Sonarr"}` is 1 and `prowlarr_app_info{name="Radarr", sync_level="fullSync", implementation="Radarr"}` is 1

#### Scenario: No connected apps
- **WHEN** no applications are configured
- **THEN** no `prowlarr_app_info` metrics are emitted

### Requirement: Connected applications total
`prowlarr_apps_total` SHALL be a gauge with the total number of connected applications.

#### Scenario: Three apps connected
- **WHEN** 3 applications are configured
- **THEN** `prowlarr_apps_total` is 3

### Requirement: Indexer failure status metric
`prowlarr_indexer_disabled{indexer}` SHALL be a gauge: 1 if the indexer is temporarily disabled (has a DisabledTill in the future), 0 otherwise. The `indexer` label SHALL use the indexer name resolved from the indexer ID.

#### Scenario: Disabled indexer
- **WHEN** indexer "NZBgeek" (ID 1) has DisabledTill in the future
- **THEN** `prowlarr_indexer_disabled{indexer="NZBgeek"}` is 1

#### Scenario: All indexers healthy
- **WHEN** no indexers have a DisabledTill in the future
- **THEN** no `prowlarr_indexer_disabled` metrics with value 1 are emitted
