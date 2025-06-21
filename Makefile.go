# ────────────────────── Makefile for Go ──────────────────────────────────────
.PHONY: build test docker-build docker-run

MOD_ROOT := demo-go-server

build:
	go build -o bin/server src/go/main.go

test:
	go test ./tests/go/...

docker-build:
	docker build -f Dockerfile.go -t demo-go-server-go:local .

docker-run:
	docker run --rm -p 50051:50051 demo-go-server-go:local
