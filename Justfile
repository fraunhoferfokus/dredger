program := `basename $PWD`
module := "build"

# Show this help message
help:
    just -l

# Build the local dredger binary
build:
    go mod tidy
    go build -o {{program}}

# Install the dredger binary in the GOPATH
install:
	go install

# Generate the source code in the target directory ./src from the OpenAPI file provided in the environment variable OPEN_API_PATH
generate:
	go run main.go generate $(OPEN_API_PATH)

# Generate the source code with all options
generate-all-flags:
	go run main.go generate $(OPEN_API_PATH) -o . -n {{module}} -l -d

# Run all tests
test:
	go test ./... -v

# Update required tools and libraries
update: download-rapidoc download-elements download-style
    go get -u
    go mod tidy

# Install additionally required tools
tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/a-h/templ/cmd/templ@latest

# Download all necessary libraries and tidy up go mod
download-deps:
    go mod download github.com/nats-io/nats.go
    go mod download github.com/rs/zerolog
    go mod download go.opentelemetry.io/otel
    go mod download github.com/nats-io/nkeys
    go mod tidy

# Download rapidoc, an OpenAPI documentation viewer
download-rapidoc:
    curl -o templates/web/js/rapidoc-min.js -L https://unpkg.com/rapidoc/dist/rapidoc-min.js

# Download elements, an OpenAPI documentation viewer (https://stoplight.io/open-source/elements)
download-elements:
    curl -o templates/web/js/elements.min.js -L https://unpkg.com/@stoplight/elements/web-components.min.js
    curl -o templates/web/css/elements.min.css -L https://unpkg.com/@stoplight/elements/styles.min.css

# Download frontend libraries
# woff Dateien müssen aus dem ZIP manuell nach fonts/ kopiert werden!
download-bootstrap:
    curl -o templates/web/css/bootstrap.min.css -L https://cdn.jsdelivr.net/npm/bootstrap@latest/dist/css/bootstrap.min.css
    curl -o templates/web/css/bootstrap-icons.min.css -L https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css
    curl -o templates/web/css/bootstrap-icons-1.11.3.zip -L https://github.com/twbs/icons/releases/download/v1.11.3/bootstrap-icons-1.11.3.zip
    curl -o templates/web/js/bootstrap.bundle.min.js -L https://cdn.jsdelivr.net/npm/bootstrap@latest/dist/js/bootstrap.bundle.min.js
download-style: download-bootstrap
    curl -o templates/web/js/htmx.min.js -L https://unpkg.com/htmx.org@2/dist/htmx.min.js
    curl -o templates/web/js/hyperscript.js -L https://unpkg.com/hyperscript.org@latest
    curl -o templates/web/js/sse.js -L https://unpkg.com/htmx-ext-sse@2/sse.js
    curl -o templates/web/js/rapidoc-min.js -L https://unpkg.com/rapidoc/dist/rapidoc-min.js
    curl -o templates/web/css/simple.min.css -L https://unpkg.com/simpledotcss/simple.min.css
    curl -o templates/web/css/pico.min.css -L https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css
    curl -o templates/web/css/pico.colors.min.css -L https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.colors.min.css
certificate:
    go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost`

# List all ToDo items of the source code
todo:
    rg -ip "to.?do:"

# Check for passwords
gitleaks:
    gitleaks --no-banner -l warn -v dir .

# push to git
push: gitleaks
    git add .
    git commit -a
    git push

# changes for git
status:
    git status -s
