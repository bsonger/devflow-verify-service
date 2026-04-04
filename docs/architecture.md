# Architecture

## Purpose

`devflow-verify-service` is the execution-fact ingress and writeback backend.
It accepts Tekton / Argo / release-step execution facts and writes them back through the approved verify path.

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

## Request Flow

```text
Observer / Controller / External callback
  -> verify router
  -> verify handler
  -> verify service
  -> Mongo-backed writeback path
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
  - manifest/release/intent writeback logic
- `pkg/store`
  - Mongo access
- `pkg/model`
  - minimal writeback-facing models

## External Dependencies

- `Gin`
- `MongoDB`
- `devflow-service-common`
- Tekton / Argo / controller callback sources

## Non-Goals

- `Project` CRUD
- `Application` CRUD
- `Configuration` CRUD
- public owner semantics for `Manifest` / `Release` / `Intent`
- Tekton / Argo active execution dispatch
