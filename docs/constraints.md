# Constraints

## Ownership

- this repo exposes only `Verify` as its public resource surface
- `Image`, `Release`, and `Intent` exist here only as internal writeback dependencies
- do not reintroduce public routers, handlers, or Swagger surfaces for unrelated resources

## Hard constraints

- writeback must target resources precisely by resource ID
- repeated writeback to terminal steps must remain idempotent
- token validation must apply to write interfaces
- verify writeback auth config must come from mounted `config.yaml`
- do not read shared tokens directly from environment variables in handlers or middleware
- Swagger must contain only `/api/v1/verify/*`

## Data rules

- `pipeline_id`, `task_name`, and `step_name` must never be empty strings
- verify payloads must continue to describe a verify-only boundary in `README.md` and `docs/*.md`

## Dependency rules

- any call to another service or external system must emit metrics, traces, and structured logs together
- do not use `release_id`, `image_id`, `intent_id`, or `external_ref` as metric labels
- those identifiers belong in logs and trace attributes only

## Non-goals

- exposing public CRUD surfaces for `Image`, `Release`, or `Intent`
- broadening repo ownership beyond verify ingress and writeback normalization
