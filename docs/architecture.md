# Architecture

## Purpose

`devflow-verify-service` is the execution-fact ingress and writeback backend.
It accepts Tekton / Argo / release-step execution facts and writes them back through the approved verify path.
It does not become the public CRUD owner for control-plane resources.

## Architecture Style

This repo uses a **layered ingress/writeback backend**:

```text
router -> api -> service -> store
                    \-> model
```

The service layer should stay focused on:
- fact normalization
- writeback rules
- minimal status/step update behavior
- verification-result persistence only

The target relational resource model is:

- `ManifestVerification` = build verification/result record
- `ReleaseVerification` = release verification/result record

## Request Flow

```text
Observer / Controller / External callback
  -> verify router
  -> verify handler
  -> verify service
  -> writeback / verification-result store
  -> HTTP response
```

## Internal Package Layout

- `cmd/main.go`
  - process entrypoint only
- `pkg/config`
  - config loading
  - runtime initialization
- `pkg/router`
  - verify route registration
  - middleware wiring
- `pkg/api`
  - verify handlers
- `pkg/service`
  - manifest/release writeback logic
  - minimal verification-result update logic
- `pkg/store`
  - verification-result persistence
- `pkg/model`
  - verification-result-facing models

## External Dependencies

- `Gin`
- PostgreSQL persistence
- `devflow-service-common`
- Tekton / Argo / controller callback sources

## Non-Goals

- `Project` CRUD
- `Application` CRUD
- `Configuration` CRUD
- public owner semantics for `Manifest` / `Release` / `Intent`
- Tekton / Argo active execution dispatch
- long-term business-source-of-truth storage for release resources

## Swagger generation

- `Dockerfile` runs `swag init -g cmd/main.go --parseDependency -o docs/generated/swagger` before building.
- Generated artifacts live in `docs/generated/swagger` and should be regenerated when handlers change.
- Export scripts copy the same folder so split repos stay consistent.
