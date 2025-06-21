package go_test

import (
	"context"
	v2 "github.com/wjb-dev/demo-go-server/pkg/proto/v1"
	"testing"
)

type stubServer struct {
	v2.UnimplementedEchoServiceServer
}

func (s *stubServer) Echo(ctx context.Context, req *v2.EchoRequest) (*v2.EchoResponse, error) {
	return &v2.EchoResponse{Message: req.Message}, nil
}

func TestEcho(t *testing.T) {
	srv := &stubServer{}
	req := &v2.EchoRequest{Message: "hello"}
	resp, err := srv.Echo(context.Background(), req)
	if err != nil {
		t.Fatalf("Echo returned error: %v", err)
	}
	if resp.Message != "hello" {
		t.Errorf("expected %q, got %q", "hello", resp.Message)
	}
}
