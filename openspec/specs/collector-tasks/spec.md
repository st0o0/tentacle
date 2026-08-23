## Purpose

Prometheus collector exposing Jellyfin scheduled task states, progress, and execution history.

## Requirements

### Requirement: Task state metric
`jellyfin_scheduled_task_state{task_name, category}` SHALL be a gauge: 0=Idle, 1=Running, 2=Cancelling.

#### Scenario: Running task
- **WHEN** task "Scan Media Library" with category "Library" is running
- **THEN** `jellyfin_scheduled_task_state{task_name="Scan Media Library", category="Library"}` is 1

#### Scenario: Idle task
- **WHEN** task "Clean Cache" is idle
- **THEN** `jellyfin_scheduled_task_state{task_name="Clean Cache", ...}` is 0

### Requirement: Task progress metric
`jellyfin_scheduled_task_progress_ratio{task_name}` SHALL be a gauge from 0 to 1 representing current task progress.

#### Scenario: Task at 75%
- **WHEN** `CurrentProgressPercentage` is 75
- **THEN** `jellyfin_scheduled_task_progress_ratio{task_name="Scan Media Library"}` is 0.75

### Requirement: Last run success metric
`jellyfin_scheduled_task_last_run_success{task_name}` SHALL be a gauge: 1 if the last execution status was "Completed", 0 otherwise.

#### Scenario: Successful last run
- **WHEN** `LastExecutionResult.Status` is "Completed"
- **THEN** `jellyfin_scheduled_task_last_run_success{task_name="Scan Media Library"}` is 1

#### Scenario: Failed last run
- **WHEN** `LastExecutionResult.Status` is "Failed"
- **THEN** `jellyfin_scheduled_task_last_run_success{task_name="Scan Media Library"}` is 0

#### Scenario: No execution history
- **WHEN** `LastExecutionResult` is null
- **THEN** no last_run_success, last_run_timestamp, or last_run_duration metrics are emitted for that task

### Requirement: Last run timestamp metric
`jellyfin_scheduled_task_last_run_timestamp_seconds{task_name}` SHALL be a gauge with the Unix timestamp of when the task last completed (EndTimeUtc).

#### Scenario: Last completed
- **WHEN** `EndTimeUtc` is "2024-01-15T10:30:00Z"
- **THEN** `jellyfin_scheduled_task_last_run_timestamp_seconds{task_name="..."}` is the Unix timestamp

### Requirement: Last run duration metric
`jellyfin_scheduled_task_last_run_duration_seconds{task_name}` SHALL be a gauge with the duration in seconds, calculated as `EndTimeUtc - StartTimeUtc`.

#### Scenario: 30-second run
- **WHEN** StartTimeUtc is 10:30:00 and EndTimeUtc is 10:30:30
- **THEN** `jellyfin_scheduled_task_last_run_duration_seconds{task_name="..."}` is 30
