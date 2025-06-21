package go_test

import (
	"context"
	"testing"

	v1 "github.com/wjb-dev/demo-go-server/proto/v1"
)

type stubServer struct {
	v1.UnimplementedEchoServiceServer
}

func (s *stubServer) Echo(ctx context.Context, req *v1.EchoRequest) (*v1.EchoResponse, error) {
	return &v1.EchoResponse{Message: req.Message}, nil
}

func TestEcho(t *testing.T) {
	srv := &stubServer{}
	req := &v1.EchoRequest{Message: "hello"}
	resp, err := srv.Echo(context.Background(), req)
	if err != nil {
		t.Fatalf("Echo returned error: %v", err)
	}
	if resp.Message != "hello" {
		t.Errorf("expected %q, got %q", "hello", resp.Message)
	}
}
