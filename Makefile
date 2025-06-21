# Universal Makefile — C++ vs Go
LANG        := go
PROJECT     := demo-go-server
GO_SRCDIR   := src/go
BIN_DIR     := bin
GO_BINARY   := $(BIN_DIR)/$(PROJECT)-go
PROTO_DIR   := proto/v1

.PHONY: generate build test docker-build docker-run clean

ifeq ($(LANG),cpp)
build:
	@mkdir -p build
	cmake -S src/cpp -B build -DCMAKE_BUILD_TYPE=Release
	cmake --build build

test:
	cd build && ctest --output-on-failure

docker-build:
	docker build -f Dockerfile.cpp -t $(PROJECT)-cpp:local .

docker-run:
	docker run --rm -p 50051:50051 $(PROJECT)-cpp:local

clean:
	rm -rf build

endif

ifeq ($(LANG),go)

# 1️⃣ Generate Go code from your .proto
generate:
	@echo "🛠️  Generating Go code from .proto…"
	protoc \
	  -I $(PROTO_DIR) \
	  --go_out=paths=source_relative:$(PROTO_DIR) \
	  --go-grpc_out=paths=source_relative:$(PROTO_DIR) \
	  $(PROTO_DIR)/service.proto

# 2️⃣ Build the Go binary (after codegen)
build: generate
	@echo "🔨  Building Go binary…"
	@mkdir -p $(dir $(BIN_DIR))
	CGO_ENABLED=0 go build -o $(BIN_DIR) ./$(GO_SRCDIR)

# 3️⃣ Run your Go unit tests
test:
	@echo "🚦  Running Go tests…"
	go test ./$(SRCDIR)/...

# 4️⃣ Build the Docker image
docker-build:
	docker build -f Dockerfile.go -t $(PROJECT)-go:local .

# 5️⃣ Run the container
docker-run:
	docker run --rm -p 50051:50051 $(PROJECT)-go:local

# 6️⃣ Clean up
clean:
	rm -rf $(BIN_DIR)
endif
