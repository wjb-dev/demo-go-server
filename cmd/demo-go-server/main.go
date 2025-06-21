package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wjb-dev/demo-go-server/internal/config"
	"github.com/wjb-dev/demo-go-server/internal/server"
)

func main() {
	// Define flags to override environment config if desired
	flagPort := flag.Int("port", 0, "gRPC server port. If set (>0), overrides GRPC_PORT env")
	flagReflection := flag.Bool("reflection", false, "Enable gRPC reflection. If set, overrides GRPC_ENABLE_REFLECTION env")
	flag.Parse()

	// Load from environment
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Override with flags if provided
	if *flagPort != 0 {
		cfg.Port = fmt.Sprintf("%d", *flagPort)
	}
	if *flagReflection {
		cfg.EnableReflection = true
	}

	// Log startup info
	log.Printf("🚀 Starting gRPC server on port %s (reflection=%v)", cfg.Port, cfg.EnableReflection)

	// Optionally: handle graceful shutdown
	// Here we run server.Run in a goroutine, and listen for signals to exit.
	// Note: for a more graceful stop (letting in-flight RPCs finish), you could refactor server.Run to return *grpc.Server.
	// For simplicity, we exit the process when signal is received.
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	// Run server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := server.Run(cfg); err != nil {
			errCh <- err
		}
	}()

	// Wait for error or shutdown signal
	select {
	case sig := <-exitCh:
		log.Printf("Received signal %v, shutting down...", sig)
		// If you need graceful stop: call server.GracefulStop() on the grpc.Server instance.
		// In this simple pattern, we just exit.
		os.Exit(0)
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	}
}
