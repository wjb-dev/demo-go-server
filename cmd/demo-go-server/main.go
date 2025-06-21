package main

import (
	"context"
	"flag"
	"fmt"
	v2 "github.com/wjb-dev/demo-go-server/pkg/proto/v1"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type echoServer struct {
	v2.UnimplementedEchoServiceServer
}

func (s *echoServer) Echo(ctx context.Context, req *v2.EchoRequest) (*v2.EchoResponse, error) {
	return &v2.EchoResponse{Message: req.Message}, nil
}

func main() {
	port := flag.Int("port", 50051, "gRPC server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	v2.RegisterEchoServiceServer(grpcServer, &echoServer{})
	reflection.Register(grpcServer)

	log.Printf("🚀  gRPC server listening at %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
