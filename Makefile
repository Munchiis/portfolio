.PHONY: build serve dev clean gen

# Generate templ templates
gen:
	templ generate

# Build the static site
build: gen
	go run cmd/portfolio/main.go --build

# Serve the static site
serve: build
	go run cmd/portfolio/main.go

# Run in development mode
dev: gen
	go run cmd/portfolio/main.go --dev

# Clean build artifacts
clean:
	rm -rf dist
	rm -rf templates/**/*_templ.go

# Install dependencies
deps:
	go install github.com/a-h/templ/cmd/templ@latest

# Generate a self-contained binary
binary: gen
	go build -o portfolio cmd/portfolio/main.go
