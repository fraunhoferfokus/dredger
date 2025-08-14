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
- `parser/` – modules that parse OpenAPI and AsyncAPI documents.
- `templates/` – embedded templates used by the generator to emit code.
- `tests/` – common test setup and fixtures.

## Go module files
- `go.mod`, `go.sum` – module dependencies for the generator itself.
- `go.work`, `go.work.sum` – workspace configuration for multi-module builds.
