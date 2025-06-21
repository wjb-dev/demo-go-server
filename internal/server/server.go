package server

import (
	"fmt"
	"net"

	"github.com/wjb-dev/demo-go-server/internal/config"
	"github.com/wjb-dev/demo-go-server/internal/handler"
	v1 "github.com/wjb-dev/demo-go-server/pkg/proto/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Run creates and starts the gRPC server according to cfg.
// It blocks serving until the server stops (e.g., fatal error or externally interrupted).
func Run(cfg *config.Config) error {
	// Listen on configured port
	addr := fmt.Sprintf(":%s", cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// Create gRPC server. You can add ServerOptions here (e.g., TLS, interceptors).
	grpcServer := grpc.NewServer()

	// Register your service handlers:
	echoSrv := handler.NewEchoServiceServer()
	v1.RegisterEchoServiceServer(grpcServer, echoSrv)

	// Register health service:
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	// Set overall status to SERVING. The empty string "" indicates the overall server.
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Register reflection if enabled:
	if cfg.EnableReflection {
		reflection.Register(grpcServer)
	}

	// Start serving (blocking call). Logging is done in main.
	return grpcServer.Serve(lis)
}
