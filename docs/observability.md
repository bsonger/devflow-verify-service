# Observability

## Shared Baseline

This repo follows the shared telemetry contract implemented in `devflow-service-common`.

- structured logs with shared runtime fields
- `devflow_http_*` ingress metrics
- standard server/client spans with service-defined business attributes
- optional diagnostics only for `pprof` and Pyroscope

## Repo-Local Focus

`devflow-verify-service` should add writeback context for:

- `manifest`
- `release`
- `intent`
- `pipeline`
- `task`

Recommended structured fields:

- `resource`
- `resource_id`
- `manifest_id`
- `release_id`
- `intent_id`
- `pipeline_id`
- `task_name`
- `result`
- `error_code`

## Async Notes

- verify writeback endpoints should preserve inbound request trace context when callers provide it
- updates triggered by external controllers should emit resource-scoped logs instead of high-cardinality metrics labels

## Profile

- `pprof` is disabled by default
- Pyroscope is disabled by default
- both are enabled only through explicit runtime configuration
