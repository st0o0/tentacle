## ADDED Requirements

### Requirement: Issue Count endpoint
`GetIssueCount()` SHALL call `GET /api/v1/issue/count` and return `IssueCount` (Total, Open, Resolved).

#### Scenario: Issue counts
- **WHEN** `/api/v1/issue/count` returns Total=10, Open=3, Resolved=7
- **THEN** `IssueCount.Total` is 10, `.Open` is 3, `.Resolved` is 7

#### Scenario: No issues
- **WHEN** `/api/v1/issue/count` returns all zeros
- **THEN** all fields are 0
