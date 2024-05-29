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
update: download-rapidoc download-style
    go get -u

# Install additionally required tools
tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/a-h/templ/cmd/templ@latest

# Download rapidoc, an OpenAPI documentation viewer
download-rapidoc:
    curl -o templates/web/js/rapidoc-min.js -L https://unpkg.com/rapidoc/dist/rapidoc-min.js

# Download frontend libraries
download-style:
    curl -o templates/web/css/bootstrap.min.css -L https://cdn.jsdelivr.net/npm/bootstrap@latest/dist/css/bootstrap.min.css
    curl -o templates/web/css/bootstrap-icons.min.css -L https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.2/font/bootstrap-icons.min.css
    curl -o templates/web/js/bootstrap.bundle.min.js -L https://cdn.jsdelivr.net/npm/bootstrap@latest/dist/js/bootstrap.bundle.min.js
    curl -o templates/web/js/htmx.min.js -L https://unpkg.com/htmx.org@latest
    curl -o templates/web/js/htmx-sse.js -L https://unpkg.com/htmx.org/dist/ext/sse.js
    curl -o templates/web/js/hyperscript.js -L https://unpkg.com/hyperscript.org@latest
    curl -o templates/web/js/sse.js -L https://unpkg.com/htmx.org/dist/ext/sse.js
    curl -o templates/web/js/rapidoc-min.js -L https://unpkg.com/rapidoc/dist/rapidoc-min.js
    curl -o templates/web/css/simple.min.css -L https://unpkg.com/simpledotcss/simple.min.css
    curl -o templates/web/css/pico.min.css -L https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css
    curl -o templates/web/css/pico.colors.min.css -L https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.colors.min.css
certificate:
    go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost`

# List all ToDo items of the source code
todo:
    rg -ip "to.?do:"

# push to git
push:
    git add .
    git commit -a
    git push

# changes for git
status:
    git status -s
