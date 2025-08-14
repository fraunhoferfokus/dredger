# File Overview

This document outlines the purpose of the main directories and files in the repository to help new contributors navigate the project.

## Top-level files
- `main.go` – entry point invoking the CLI to run the generator.
- `Justfile` – handy recipes for building, testing and generating code.
- `Features.md` – high level list of supported generator features.
- `Testing-BDD.md` – notes on behaviour-driven tests.
- `version` – current release number embedded into the binaries.

## Core directories
- `cli/` – command-line interface definitions and flag handling.
- `cmd/` – wrappers for executing external tooling during generation.
- `core/` – shared helpers such as embedded version information.
- `dredger/` – prebuilt binary of the generator.
- `examples/` – sample OpenAPI and AsyncAPI specs used for experiments.
- `fileUtils/` – utilities for working with files and directories.
- `generator/` – logic that creates source code from parsed specifications.
  - `asyncapi/` – helpers for processing AsyncAPI specs and generating messaging code.
  - `features/` – sample BDD feature files used to validate generation behaviour.
  - `tests/` – unit tests for generator utilities.
  - `bddOpen.go` – generates Godog test stubs from feature files.
  - `channelCommon.go` – shared helpers for AsyncAPI channel processing.
  - `config.go` – writes configuration sources such as `.env` and `config.go`.
  - `database.go` – optional SQLite scaffolding.
  - `dockerfile.go` – emits Dockerfile and image manifest.
  - `entitiesAsync.go` – derives Go types from AsyncAPI message schemas.
  - `entitiesOpen.go` – builds entity structs and validators from OpenAPI schemas.
  - `frontendAsync.go` – sets up minimal frontend for messaging projects.
  - `frontendOpen.go` – sets up minimal frontend for HTTP projects.
  - `generatorAsync.go` – orchestrates AsyncAPI project generation.
  - `generatorOpen.go` – orchestrates OpenAPI project generation.
  - `godog_test.go` – integration tests for feature-based generation.
  - `handlerOpen.go` – creates HTTP handler stubs.
  - `info.go` – renders runtime info endpoint.
  - `justfile.go` – writes project `Justfile`.
  - `lifecycleOpen.go` – adds liveness and readiness endpoints.
  - `logger.go` – configures zerolog logging.
  - `oapi.go` – high-level OpenAPI generation entry.
  - `policy.go` – scaffolds policy and authorization hooks.
  - `readme.go` – generates project `README.md`.
  - `serverAsync.go` – NATS server and subscriber setup for AsyncAPI.
  - `templates.go` – loads embedded code generation templates.
  - `tracing.go` – wires OpenTelemetry tracing and metrics.
  - `types.go` – CLI flag definitions used by the generator.
  - `utils.go` – assorted helper functions.
  - `validation.go` – validation helpers for entity data.
- `parser/` – modules that parse OpenAPI and AsyncAPI documents.
- `templates/` – embedded templates used by the generator to emit code.
- `tests/` – common test setup and fixtures.

## Go module files
- `go.mod`, `go.sum` – module dependencies for the generator itself.
- `go.work`, `go.work.sum` – workspace configuration for multi-module builds.
