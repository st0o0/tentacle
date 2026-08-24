## ADDED Requirements

### Requirement: Issues total metric
`seerr_issues_total` SHALL be a gauge with the total number of reported issues.

#### Scenario: Issues exist
- **WHEN** Seerr has 10 total issues
- **THEN** `seerr_issues_total` is 10

#### Scenario: No issues
- **WHEN** Seerr has 0 issues
- **THEN** `seerr_issues_total` is 0

### Requirement: Issues by status metric
`seerr_issues_by_status_total{status}` SHALL be a gauge with issue counts by status (open, resolved).

#### Scenario: Mixed issue statuses
- **WHEN** Seerr has 3 open and 7 resolved issues
- **THEN** `seerr_issues_by_status_total{status="open"}` is 3 and `{status="resolved"}` is 7

#### Scenario: All resolved
- **WHEN** all 10 issues are resolved
- **THEN** `seerr_issues_by_status_total{status="open"}` is 0 and `{status="resolved"}` is 10

#### Scenario: Issues API failure
- **WHEN** `GET /api/v1/issue/count` fails
- **THEN** no issue metrics are emitted and `seerr_scrape_success{collector="issues"}` is 0
