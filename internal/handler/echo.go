package handler

import (
	"context"

	v1 "github.com/wjb-dev/demo-go-server/pkg/proto/v1"
)

// EchoServiceServer implements the EchoService defined in proto.
type EchoServiceServer struct {
	v1.UnimplementedEchoServiceServer
	// You can add dependencies here, e.g., logger, metrics, DB client, etc.
}

// NewEchoServiceServer constructs a new EchoServiceServer. Add dependencies as parameters if needed.
func NewEchoServiceServer() *EchoServiceServer {
	return &EchoServiceServer{}
}

// Echo simply echoes back the incoming message.
func (s *EchoServiceServer) Echo(ctx context.Context, req *v1.EchoRequest) (*v1.EchoResponse, error) {
	// Business logic: here just echo the message
	return &v1.EchoResponse{Message: req.Message}, nil
}
