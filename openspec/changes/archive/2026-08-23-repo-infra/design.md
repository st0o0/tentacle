## Approach

Direct adaptation of the ran repository infrastructure. Every file is copied and search-replaced for tentacle-specific values. No structural changes to the workflow patterns — they're proven in production on ran.

## Key Decisions

### No e2e job
Tentacle requires a running Jellyfin instance with an API token. Unlike ran (which can self-test with netcat against its own listeners), we can't spin up a Jellyfin server in CI cheaply. The smoke test validates the binary builds and runs (`tentacle version`).

### Smoke test strategy
Ran's smoke test runs the binary and expects exit code 1 (no traps configured). Tentacle's equivalent: run `tentacle version` which prints the version and exits 0 — this proves the binary is valid without needing TENTACLE_JELLYFIN_ADDRESS.

### Initial version 0.1.0
Start at 0.1.0 in the release-please manifest. First release will be 0.1.0 or 0.2.0 depending on whether we have a feat commit.

### Same commitlint rules
Identical conventional commit types. The `header-max-length` is a warning (level 1), not an error. Dependabot commits are ignored.

## File Mapping

```
ran                              → tentacle
─────────────────────────────────────────────────
.github/workflows/ci.yml        → adapted (no e2e, smoke = version)
.github/workflows/release.yml   → adapted (image name, description)
.github/workflows/security.yml  → adapted (image name)
.github/workflows/commitlint.yml → identical
.github/workflows/dev-build.yml → adapted (image name, description)
release-please-config.json      → identical
.release-please-manifest.json   → version 0.1.0
commitlint.config.mjs           → identical
package.json                    → identical
.hadolint.yaml                  → identical
LICENSE.md                      → identical
CODE_OF_CONDUCT.md              → identical
CONTRIBUTING.md                 → adapted (project name, commands)
SECURITY.md                     → adapted (image name, description)
```
