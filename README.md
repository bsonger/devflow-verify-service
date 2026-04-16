# DevFlow Verify Service

`devflow-verify-service` is the backend execution-fact ingress for verify writeback.

## Role

- expose `/api/v1/verify/*`
- accept external execution facts
- normalize verify/writeback payloads
- write execution facts back to control-plane resources through the approved path

## Key Commands

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
- Swagger UI: `/swagger/index.html`
- Staging Swagger UI: `/api/v1/verify/swagger/index.html`

## Key Docs

- `docs/README.md`
- `scripts/README.md`
- `docs/architecture.md`
- `docs/constraints.md`
- `docs/observability.md`
- `docs/api-spec.md`
- `docs/resources/README.md`
- `docs/generated/swagger/swagger.yaml`
