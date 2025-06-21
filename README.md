# Demo Go Server

A cross-language gRPC microservice.

Local mode (docker - recommended):
make run

Local mode (MacOs):

brew install go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$(go env GOPATH)/bin:$PATH"
make generate-local


