# Build stage for generating templ files
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install tools we need
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go.mod first for better caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ files
RUN templ generate

# Build the binary
RUN go build -o portfolio ./cmd/portfolio

# Development stage
FROM golang:1.24-alpine AS development

WORKDIR /app

# Install templ for development hot reload
RUN go install github.com/a-h/templ/cmd/templ@latest

# We'll mount the source code as a volume
EXPOSE 3000

CMD ["go", "run", "./cmd/portfolio/main.go", "--dev", "--port", "3000"]

# Production stage - lightweight alpine
FROM alpine:latest AS production

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/portfolio /app/

# Copy static files
COPY static/ /app/static/

# Create dist directory
RUN mkdir -p /app/dist

# Build the static site when container starts
# Then serve the static files
EXPOSE 8080

CMD ["/app/portfolio", "--build", "--dir", "/app/dist", "--port", "8080"]
