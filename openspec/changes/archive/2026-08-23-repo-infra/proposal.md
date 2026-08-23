## Summary

Add CI/CD pipelines, release automation, and community files to the tentacle repository, mirroring the established pattern from the ran repository. This brings tentacle to the same operational maturity: every PR is validated, releases are automated, images are signed, and the repo has standard open-source governance files.

## Motivation

The tentacle repo currently has code and a Dockerfile but no CI, no automated releases, and no community files. The ran repo has a proven setup that we can adapt directly — same Go toolchain, same Docker build pattern, same conventions.

## Scope

### CI/CD Workflows (`.github/workflows/`)
- **ci.yml** — PR validation: unit tests, golangci-lint + hadolint, Docker build, smoke test
- **release.yml** — release-please on main → multi-arch GHCR push + cosign signing
- **security.yml** — weekly Trivy scan + on PR when Dockerfile/deps change
- **commitlint.yml** — conventional commits enforcement on PRs
- **dev-build.yml** — label-triggered dev image push to GHCR with environment approval gate

### Config Files
- **release-please-config.json** — simple release type, pre-major bump rules
- **.release-please-manifest.json** — initial version `0.1.0`
- **commitlint.config.mjs** — conventional commit types + rules
- **package.json** — commitlint devDependency only
- **.hadolint.yaml** — Dockerfile linter (failure-threshold: warning, ignore DL3018)

### Community Files
- **LICENSE.md** — MIT
- **CODE_OF_CONDUCT.md** — Contributor Covenant 2.1
- **CONTRIBUTING.md** — dev setup, PR guidelines, conventions
- **SECURITY.md** — supported versions, vulnerability reporting

## Non-goals

- No e2e test suite (tentacle needs a Jellyfin instance — not practical in CI without a mock)
- No changes to application code
- No README rewrite (separate concern)

## Adaptations from ran

All files are copied from ran and adapted for tentacle:
- Image name: `ghcr.io/st0o0/tentacle`
- Binary name: `tentacle`
- Description: "Prometheus exporter for Jellyfin media servers"
- Default port: 9594
- No e2e job in ci.yml (no test suite yet)
- Smoke test: `tentacle version` instead of exit-code check (tentacle needs env vars to run)
