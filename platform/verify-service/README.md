# Verify Service Platform Notes

## Purpose

This file is the repo-local runtime note for `devflow-verify-service`.
For public API shape, ownership, and request payload details, prefer:
- `../README.md`
- `../docs/`
- `../docs/resources/verify.md`

## Runtime entrypoints

- process entry: `cmd/main.go`
- shared bootstrap: `../devflow-service-common/bootstrap`
- router root: `pkg/router/router.go`

## Main local code paths

- verify routes: `pkg/router/verify.go`
- verify handler: `pkg/api/verify.go`
- manifest writeback: `pkg/service/manifest.go`
- release writeback: `pkg/service/release.go`
- intent writeback: `pkg/service/intent.go`

## Platform dependencies

- shared middleware / startup / observability moved to `devflow-service-common`
- controller integration details remain in `controller-integration.md`

## Service identity

- OTel `service.name`: `verify-service`
- verify-specific token header: `X-Devflow-Verify-Token`
