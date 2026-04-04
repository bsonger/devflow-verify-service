# Devflow Verify Service

`devflow-verify-service` is the backend execution-fact ingress for verify writeback.

## Backend Role

- expose `/api/v1/verify/*`
- accept external execution facts
- normalize verify/writeback payloads
- write execution facts back to control-plane resources through the approved path

## Backend Architecture

This repo uses a **layered ingress/writeback backend**:

```text
cmd
 -> config
 -> router
 -> api
 -> service
 -> store
 -> model
```

### Package responsibilities

- `cmd/`: service startup
- `pkg/config`: config loading and runtime init
- `pkg/router`: verify route and middleware wiring
- `pkg/api`: verify handlers
- `pkg/service`: normalization and writeback logic
- `pkg/store`: Mongo access for writeback
- `pkg/model`: minimal models needed for writeback

## Non-Goals

- no public CRUD for `Project`
- no public CRUD for `Application`
- no public CRUD for `Configuration`
- no public owner semantics for `Manifest` / `Release` / `Intent`
- no active execution orchestration

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
