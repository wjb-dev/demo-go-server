# ──────────────────── Stage 1: Builder ────────────────────────────────────────
FROM golang:1.24.4-alpine AS builder

# Install protoc and Git
RUN apk add --no-cache git protobuf protoc

# Install the Go protoc plugins into $GOPATH/bin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure the plugins are on PATH
ENV PATH="${PATH}:$(go env GOPATH)/bin"

WORKDIR /cmd

# Copy go.mod & go.sum, download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy in your .proto and Go sources
COPY proto/ proto/
COPY cmd/demo-go-server/ cmd/demo-go-server/

# Generate Go code from your proto
RUN protoc \
-I pkg/proto/v1 \
--go_out=paths=source_relative:pkg/proto/v1 \
--go-grpc_out=paths=source_relative:pkg/proto/v1 \
pkg/proto/v1/service.proto

# Build static Go binary
RUN CGO_ENABLED=0 go build -o /app/demo-go-server cmd/demo-go-server/main.go


# ──────────────────── Stage 2: Runtime ────────────────────────────────────────
FROM alpine:latest

# Only CA certs (for TLS) & your binary
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/demo-go-server .

EXPOSE 50051
ENTRYPOINT ["./demo-go-server"]
