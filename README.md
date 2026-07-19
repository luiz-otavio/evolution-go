# evolution-go

Go SDK for integrating with Evolution API.

This repository organizes the project modules by responsibility (API client, tests, helper scripts), with a focus on:

- ease of use for WhatsApp integrations
- strong typing for HTTP contracts
- contract and integration tests
- incremental evolution based on Evolution API upstream contracts

## Overview

The goal of this project is to provide a Go layer for consuming Evolution API endpoints with local payload validation, typed responses, and test mocks.

## Workspace modules

### `api/`

Main SDK module.

Includes:

- HTTP client and configuration (`NewEvolutionClient`)
- domain services (instance, message, chat, group, settings, etc.)
- request/response models
- mocks for unit tests
- test suite (contract + integration)

Full module documentation: [api/README.md](api/README.md)

### `scripts/`

Helper scripts for development workflows.

Currently:

- `find-affected-services.sh`: helps identify services affected by changes

## Workspace structure

- `go.work`: defines active Go modules in the workspace
- `api/go.mod`: dependencies and Go version for the `api` module

## Quick start

Prerequisites:

- Go 1.26+

Install dependencies for the API module:

```bash
cd api
go mod tidy
```

Run contract/unit tests:

```bash
go test ./... -count=1
```

## Integration tests

Integration tests in the `api` module use Docker (testcontainers) and environment variables.

See:

- [api/README.md](api/README.md)
- [api/.env](api/.env)

Typical run:

```bash
cd api
go test ./... -run Integration -v -count=1
```

## References

- Evolution API documentation: https://doc.evolution-api.com/v2/
- Official Evolution repository: https://github.com/evolution-foundation/evolution-api