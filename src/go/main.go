package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	v1 "github.com/wjb-dev/demo-go-server/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type echoServer struct {
	v1.UnimplementedEchoServiceServer
}

func (s *echoServer) Echo(ctx context.Context, req *v1.EchoRequest) (*v1.EchoResponse, error) {
	return &v1.EchoResponse{Message: req.Message}, nil
}

func main() {
	port := flag.Int("port", 50051, "gRPC server port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	v1.RegisterEchoServiceServer(grpcServer, &echoServer{})
	reflection.Register(grpcServer)

	log.Printf("🚀  gRPC server listening at %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
