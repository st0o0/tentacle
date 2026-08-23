## Tasks

- [x] Task 1: CI workflow — Create `.github/workflows/ci.yml` with jobs: unit (go test -race), lint (golangci-lint + hadolint), build (Docker build + smoke test with `tentacle version`). No e2e job.
- [x] Task 2: Release workflow — Create `.github/workflows/release.yml` with release-please + multi-arch Docker push to `ghcr.io/st0o0/tentacle` + cosign signing. Adapt image name, description, and OCI annotations.
- [x] Task 3: Security workflow — Create `.github/workflows/security.yml` with Trivy scan. Adapt image tag to `tentacle:scan`.
- [x] Task 4: Commitlint workflow and config — Create `.github/workflows/commitlint.yml`, `commitlint.config.mjs`, and `package.json` — identical to ran.
- [x] Task 5: Dev-build workflow — Create `.github/workflows/dev-build.yml`. Adapt image name to `ghcr.io/st0o0/tentacle`, description to Prometheus exporter.
- [x] Task 6: Release-please config — Create `release-please-config.json` (identical to ran) and `.release-please-manifest.json` with initial version `0.1.0`.
- [x] Task 7: Hadolint config — Create `.hadolint.yaml` — identical to ran (failure-threshold: warning, ignore DL3018).
- [x] Task 8: Community files — Create `LICENSE.md` (MIT), `CODE_OF_CONDUCT.md` (Contributor Covenant 2.1), `CONTRIBUTING.md` (adapted for tentacle), `SECURITY.md` (adapted for tentacle image name and description).
