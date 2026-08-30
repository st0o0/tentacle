# Changelog

## [0.1.6](https://github.com/st0o0/tentacle/compare/v0.1.5...v0.1.6) (2026-08-30)


### Bug Fixes

* read Sonarr stats from nested statistics object and Jellyfin size from MediaSources ([c466f73](https://github.com/st0o0/tentacle/commit/c466f73fecd40f67d7bfe60c04ea6596500dd5f2))

## [0.1.5](https://github.com/st0o0/tentacle/compare/v0.1.4...v0.1.5) (2026-08-26)


### Bug Fixes

* prevent empty scrapes by doubling handler timeout and configurable batch size ([1d78e7f](https://github.com/st0o0/tentacle/commit/1d78e7f820a87054762896ec39f0d0a4dea21629))

## [0.1.4](https://github.com/st0o0/tentacle/compare/v0.1.3...v0.1.4) (2026-08-26)


### Features

* parallel collectors, scrape timeout, and request logging ([33c4278](https://github.com/st0o0/tentacle/commit/33c4278a3958baf66da29f54b2283b9e3468e090))

## [0.1.3](https://github.com/st0o0/tentacle/compare/v0.1.2...v0.1.3) (2026-08-26)


### Features

* log transient errors (503, timeouts) as warn instead of error ([3ebcd95](https://github.com/st0o0/tentacle/commit/3ebcd95df3c22796e3d28409d8f4aec4dd515c41))

## [0.1.2](https://github.com/st0o0/tentacle/compare/v0.1.1...v0.1.2) (2026-08-25)


### Bug Fixes

* audiobookshelf users unmarshal and prowlarr API version ([3ccb753](https://github.com/st0o0/tentacle/commit/3ccb7532ce47b40e6c5b84012807797dbde6f063))
* skip unsupported jellyfin libraries and add service name to error logs ([361ff95](https://github.com/st0o0/tentacle/commit/361ff95fb84f576f3ccc753e4be1fec6f0e4d0e0))

## [0.1.1](https://github.com/st0o0/tentacle/compare/v0.1.0...v0.1.1) (2026-08-25)


### Bug Fixes

* use meta-collector pattern to fix duplicate descriptor panic ([140e0e8](https://github.com/st0o0/tentacle/commit/140e0e8e1e46fc486eb02ea915468927c8ea1fd6))

## [0.1.0](https://github.com/st0o0/tentacle/compare/v0.1.0...v0.1.0) (2026-08-24)


### Features

* add extras, apps, issues collectors and enhance existing metrics ([347e852](https://github.com/st0o0/tentacle/commit/347e85285e0ee0b991532b08adf503e3e888013f))
* unified multi-service prometheus exporter ([3b0c54e](https://github.com/st0o0/tentacle/commit/3b0c54ee5acf65146fd61ee89aa9f47cedfe8528))


### Bug Fixes

* handle errcheck lint errors across test files and server.go ([f3ef33b](https://github.com/st0o0/tentacle/commit/f3ef33bca5df5829ead5f1475d934503919dfa28))


### Documentation

* update README with complete metrics reference and CLI docs ([53e4285](https://github.com/st0o0/tentacle/commit/53e42855718262a295624192930189f2caf3efab))
