package integration

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/wjb-dev/demo-go-server/internal/handler"
	v1 "github.com/wjb-dev/demo-go-server/pkg/proto/v1"
)

const bufSize = 1024 * 1024

// dialer starts an in-memory gRPC server using bufconn and returns a client connection.
// It also returns a cleanup function to stop the server and close resources.
func dialer() (*grpc.ClientConn, func(), error) {
	// Create a bufconn listener
	lis := bufconn.Listen(bufSize)

	// Create gRPC server and register our EchoService
	grpcServer := grpc.NewServer()
	echoSrv := handler.NewEchoServiceServer()
	v1.RegisterEchoServiceServer(grpcServer, echoSrv)

	// Start serving in a goroutine
	go func() {
		// Serve will block until lis is closed or server is stopped
		if err := grpcServer.Serve(lis); err != nil {
			// In tests, panic so we know if Serve fails unexpectedly
			panic("bufconn Serve error: " + err.Error())
		}
	}()

	// Create a dialer function for bufconn
	ctx := context.Background()
	dialCtx := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	// Dial a client connection
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialCtx), grpc.WithInsecure())
	if err != nil {
		// If dialing fails, stop server and listener
		lis.Close()
		grpcServer.Stop()
		return nil, nil, err
	}

	// cleanup closes connection and stops the server
	cleanup := func() {
		conn.Close()
		grpcServer.GracefulStop() // or Stop()
		lis.Close()
	}

	return conn, cleanup, nil
}

func TestEchoIntegration(t *testing.T) {
	conn, cleanup, err := dialer()
	if err != nil {
		t.Fatalf("failed to dial bufconn: %v", err)
	}
	defer cleanup()

	client := v1.NewEchoServiceClient(conn)

	// Set a deadline to avoid hangs
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Perform the Echo RPC
	req := &v1.EchoRequest{Message: "hello integration"}
	resp, err := client.Echo(ctx, req)
	if err != nil {
		t.Fatalf("Echo RPC failed: %v", err)
	}
	if resp.Message != req.Message {
		t.Errorf("expected message %q, got %q", req.Message, resp.Message)
	}
}
