# Universal Makefile — C++ vs Go

# Project Configuration
LANG                   := go
PROJECT                := demo-go-server
GO_TESTDIR             := test/integration
BIN_DIR                := bin
PROTO_DIR              := pkg/proto/v1
PLATFORM               := $(shell uname -m)

.PHONY: install build test docker-build docker-run clean generate-proto clean-proto clean-all

# C++ targets
ifeq ($(LANG),cpp)

build:
	@echo "🔨 Building C++ binary…"
	@mkdir -p build
	cmake -S src/cpp -B build -DCMAKE_BUILD_TYPE=Release
	cmake --build build

test:
	@echo "✅ Running C++ tests…"
	cd build && ctest --output-on-failure

docker-build:
	@echo "🐳 Building C++ Docker image…"
	docker build -f Dockerfile.cpp -t $(PROJECT)-cpp:local .

docker-run: docker-build
	@echo "🚀 Running C++ Docker container…"
	docker run --rm -p 50051:50051 $(PROJECT)-cpp:local

clean:
	rm -rf build

endif

ifeq ($(LANG),go)

install:
	@echo "🛠️  Installing Go dependencies…"

ifeq ($(PLATFORM),arm64)
	@if ! command -v go >/dev/null 2>&1; then \
		echo "🍎 Apple Silicon detected, installing Go via Homebrew..."; \
		brew install go; \
	else \
		echo "✅ Go is already installed."; \
	fi
endif

	@for pkg in \
		google.golang.org/protobuf/cmd/protoc-gen-go@latest \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; do \
		binary=$$(basename $$pkg | cut -d'@' -f1); \
		if ! command -v $$binary >/dev/null 2>&1; then \
			echo "⬇️  Installing $$binary..."; \
			go install $$pkg; \
		else \
			echo "✅ $$binary already installed."; \
		fi; \
	done

generate-proto: install
	@echo "⚙️  Generating Go code from .proto…"
	PATH="$(shell go env GOPATH)/bin:$(PATH)" protoc \
		-I $(PROTO_DIR) \
		--go_out=paths=source_relative:$(PROTO_DIR) \
		--go-grpc_out=paths=source_relative:$(PROTO_DIR) \
		$(PROTO_DIR)/service.proto

test:
	@echo "✅ Running Go tests…"
	cd $(GO_TESTDIR) && go test ./...

build:
	@echo "🔨 Building Go Docker image…"
	docker build -f Dockerfile.go -t $(PROJECT)-go:local .

run: build
	@echo "🚀 Running Go Docker container…"
	docker run --rm -p 50051:50051 $(PROJECT)-go:local

clean:
	@echo "🧹 Cleaning build artifacts…"
	rm -rf $(BIN_DIR)

clean-proto:
	@echo "🧹 Cleaning generated proto files…"
	rm -rf $(PROTO_DIR)/*.pb.go

clean-all: clean clean-proto

endif