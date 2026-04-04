# Repository Guidelines

## Boundary

- This repository is `devflow-verify-service`.
- Public surface is `Verify` only.
- Internal writeback support may keep minimal `Manifest` / `Release` / `Intent` model and service code.
- Do not reintroduce `Project`, `Application`, `Configuration`, `Manifest`, `Release`, or `Intent` CRUD routes.

## Structure

- `cmd/main.go` uses shared bootstrap from `../devflow-service-common`.
- `pkg/api/` contains verify handlers only.
- `pkg/service/` contains verify writeback logic only.
- `pkg/router/` contains verify-only routes and middleware assembly.
- `pkg/config/` initializes config, observability, Mongo, and local store state.
- `docs/` contains the repository-level architecture, API, constraints, observability, and harness docs.

## Required Rules

- Any outbound service or external call must emit `metrics + trace + structured log`.
- Do not add high-cardinality business IDs to metrics labels.
- Default harness is `Planner -> Generator -> Evaluator`.
- When the runtime supports delegation, the harness must spawn those roles as separate sub-agents.
- Non-trivial work should use a run directory under `agents/runs/`.

## Doc And API Hygiene

- Regenerate Swagger after route or handler changes.
- Keep `README.md`, `AGENTS.md`, `agents/protocols/startup.md`, and `docs/*.md` aligned with the actual boundary.
- Do not reintroduce dead router/service/model/bootstrap files.
