# DevFlow Verify Service

`devflow-verify-service` is the backend execution-fact ingress for verify writeback.

## Backend Role

- expose `/api/v1/verify/*`
- accept external execution facts
- normalize verify/writeback payloads
- write execution facts back to control-plane resources through the approved path

## Local Run

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
- Swagger UI: `/swagger/index.html`

## Key Docs

- `docs/architecture.md`
- `docs/api-spec.md`
- `docs/constraints.md`
- `docs/resources/README.md`
- `docs/generated/swagger/swagger.yaml`
