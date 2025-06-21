# ──────────────────── Stage 1: Builder ────────────────────────────────────────
FROM golang:1.21-alpine AS builder

# Install protoc and Git
RUN apk add --no-cache git protobuf protoc

# Install the Go protoc plugins into $GOPATH/bin
RUN go install \
google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1 \
google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Ensure the plugins are on PATH
ENV PATH="${PATH}:$(go env GOPATH)/bin"

WORKDIR /src

# Copy go.mod & go.sum, download deps
COPY src/go/go.mod src/go/go.sum ./
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1 \
google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
ENV PATH="$PATH:$(go env GOPATH)/bin"


# Copy in your .proto and Go sources
COPY proto/ proto/
COPY src/go/ src/go/

# Generate Go code from your proto
RUN protoc \
-I proto/v1 \
--go_out=paths=source_relative:src/go \
--go-grpc_out=paths=source_relative:src/go \
proto/v1/service.proto

# Build static Go binary
RUN CGO_ENABLED=0 go build -o /app/demo-go-server src/go/main.go


# ──────────────────── Stage 2: Runtime ────────────────────────────────────────
FROM alpine:latest

# Only CA certs (for TLS) & your binary
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/demo-go-server .

EXPOSE 50051
ENTRYPOINT ["./demo-go-server"]
