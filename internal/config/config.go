package config

import (
	"os"
)

// Config holds server configuration.
type Config struct {
	Port             string // e.g., "50051"
	EnableReflection bool   // whether to register gRPC reflection
	// TODO: add more fields, e.g. TLS settings, log level, etc.
}

// LoadConfig reads configuration from environment variables (with defaults).
func LoadConfig() (*Config, error) {
	// Port
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	// Reflection
	enableReflection := false
	if os.Getenv("GRPC_ENABLE_REFLECTION") == "true" {
		enableReflection = true
	}
	// Add more env-based settings here as needed.
	return &Config{
		Port:             port,
		EnableReflection: enableReflection,
	}, nil
}

// (Optional) You can add helper to override config from flags in main, if desired.
