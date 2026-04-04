# Devflow Verify Service

This repository was exported from the `bsonger/devflow` monorepo.

GitHub target:

- `git@github.com:bsonger/devflow-verify-service.git`

Go module:

- `github.com/bsonger/devflow-verify-service`

Current scope:

- service entrypoint from `platform/verify-service/cmd/main.go`
- shared bootstrap from `platform/shared/bootstrap`
- current shared domain/runtime packages from `pkg/`

Notes:

- This is a first-stage split repo.
- Shared packages are still copied from the monorepo so the service can compile independently.
- A later cleanup phase can move stable shared pieces into `devflow-common` or another shared module.
