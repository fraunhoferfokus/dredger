# dredger

A generator for REST APIs from a given <a href="https://www.openapis.org/">OpenAPI 3</a> Specification file in either JSON or YAML format. The HTTP server uses Go's <a href="https://echo.labstack.com/">Echo</a> HTTP server as base.

This is a fork of https://github.com/MVA-OpenApi/go-open-api-generator.

# Purpose

We aim to make the life of Golang REST API developers (or non technical users) easier by creating a tool which takes an OpenAPI 3 Specification file as input and generates a basic project structure from it so that the developers can focus on the business logic. But this code could also be used by other code generators (low code) to add code using their models to create application specific micro services.

The code generation uses Go text templates to generate the code. Therefore, the code can be easily modified and extended.

Basically, the generator focuses creating the core for the API handling with their endpoints and handlers. There is basic support for integration of HTML pages and frontend libraries. It also supports typical security functions for authentication and authorisation, monitoring, logging, testing and integration with the Kubernetes eco system.

Details about the supported features could be found [here](./Features.md).

# Prerequisites

Golang (You can find an installation guide for Golang <a href="https://go.dev/">here</a>).

Godog (Only for BDD testing. You can find an installation guide here [godog](https://github.com/cucumber/godog)).

Prerequisite for HTTP/2 is a TLS connection, to generate a quick localhost certificate use either openssl or `go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost`.

[gitleaks](https://github.com/gitleaks/gitleaks)

[just](https://github.com/casey/just)

# Usage

_dredger_ is a command line tool:

    Available Commands:

-   completion Generate the autocompletion script for the specified shell
-   generate Create server and client API code from OpenApi Spec
-   generate-bdd Create BDD test file from the feature file
-   help Help about any command

Flags:
-h, --help help for dredger
Generates a REST API template from a given OpenAPI Specification file. Let's take one of the example files that are already in the project. For the sake of convenience we're going to be using `stores.yaml`.

You can check how the file looks like <a href="./examples/stores.yaml">here</a></br>

-   Step 1: We navigate to the repository folder
-   Step 2: Run the command `go run main.go generate ./examples/stores.yaml -o ./build -n build -f -D`. A description of the flags can be found below.
-   Step 3: We can now navigate to the output folder (in this case `build`) and run `go run main.go` to launch the REST API.

Generation flags:

-   `-f`               Add frontend code.
-   `-o [Output path]` Specifies the output path for the generated REST API.
-   `-n [Module name]` Specifies the go module name.
-   `-D`               Generates boilerplate code for a basic SQLite database.

For typical tasks you can use the [just](https://just.systems/man/en/) recipes:

    build              # Build the local dredger binary
    download-elements  # Download elements, an OpenAPI documentation viewer
    download-rapidoc   # Download rapidoc, an OpenAPI documentation viewer
    download-style     # Download frontend libraries
    generate           # Generate the source code in the target directory ./src from the OpenAPI file provided in the environment variable OPEN_API_PATH
    generate-all-flags # Generate the source code with all options
    help               # Show this help message
    install            # Install the dredger binary in the GOPATH
    test               # Run all tests
    tools              # Install additionally required tools
    update             # Update required tools and libraries

-   `just generate OPEN_API_PATH=path/to/open-api-file`. This command will generate the minimum project structure (no optional flags are set). The parameter `OPEN_API_PATH` is required.
-   `just generate-all-flags OPEN_API_PATH=path/to/open-api-file MODULE_NAME=module-name`. This command will generate the maximum project structure (all optional flags are set). The parameter `OPEN_API_PATH` is required.
-   `just build OUTPUT_NAME=executable-name`. This command will build an executable which can be used by the developer outside of the project repository.
-   `just test`. This command runs the unit tests for the generator.

# Examples

You can find a few OpenAPI 3 Specification file examples [here](./examples). There is also a minimal [OpenAPI.yaml](./examples/OpenAPI.yaml.min-example) file as starting point for your service.

## AsyncAPI Usage

AsyncAPI specifications are supported as well. The repository now includes a
small schema file at `examples/schemas/asyncapiv3Schema.json` used for basic
validation. To generate code from an AsyncAPI v3 file, run:

```bash
go run main.go generate ./examples/simple/asyncapiv3.json -o ./build-asyncapi -n async-service
```

To generate both an OpenAPI service and an additional AsyncAPI service in one
go, pass all spec files as positional arguments. The CLI automatically detects
whether a file is an OpenAPI or AsyncAPI specification:

```bash
go run main.go generate ./examples/stores/stores.yaml \
  ./examples/simple/asyncapiv3.json -o ./build-both -n multi-service
```
When copying the command ensure that the line break uses a `\` at the end of the
first line **without any trailing spaces**, otherwise an extra argument may be
passed to the CLI.

You can also pass several specs at once, mixing OpenAPI and AsyncAPI files:

```bash
go run main.go generate ./spec1.yaml ./spec2.yaml ./async1.json ./async2.json \
  -o ./build-all -n multi-service
```

If the schema file cannot be found, validation is skipped and the code is still
generated.

# Contributions

The origin of this project was made by 6 students (A. Uluc, A. Munteau, O. Rosenblatt, J. Wilke, C. Szramek, F. Yzeiri) of the TU Berlin as part of the module "Moderne Verteilte Anwendungen Programmierpraktikum" when studying B.Sc Computer Science and could be found at https://github.com/MVA-OpenApi/go-open-api-generator.

Further contributors: J. Gottschick
