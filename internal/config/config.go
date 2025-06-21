// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// ServerConfig holds gRPC server settings.
type ServerConfig struct {
	Host             string `yaml:"host"`              // e.g., "0.0.0.0"
	Port             int    `yaml:"port"`              // e.g., 50051
	EnableReflection bool   `yaml:"enable_reflection"` // e.g., true in dev
	MaxRecvMsgSize   int    `yaml:"max_recv_msg_size"` // e.g., 4194304
	MaxSendMsgSize   int    `yaml:"max_send_msg_size"` // e.g., 4194304
}

// TLSConfig holds TLS/mTLS settings.
type TLSConfig struct {
	Enabled      bool   `yaml:"enabled"`        // true to enable TLS
	CertFile     string `yaml:"cert_file"`      // path to server cert
	KeyFile      string `yaml:"key_file"`       // path to server key
	ClientCAFile string `yaml:"client_ca_file"` // path to CA for client cert verification
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level  string `yaml:"level"`  // e.g., "debug", "info"
	Format string `yaml:"format"` // e.g., "text" or "json"
	Output string `yaml:"output"` // e.g., "stdout" or file path
}

// MetricsConfig holds Prometheus (or other) metrics settings.
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"` // true to run metrics server
	Host    string `yaml:"host"`    // e.g., "0.0.0.0"
	Port    int    `yaml:"port"`    // e.g., 9090
	Path    string `yaml:"path"`    // e.g., "/metrics"
}

// TracingConfig holds tracing settings (e.g., OpenTelemetry).
type TracingConfig struct {
	Enabled     bool   `yaml:"enabled"`      // true to enable tracing
	Endpoint    string `yaml:"endpoint"`     // e.g., OTLP collector URL
	ServiceName string `yaml:"service_name"` // e.g., "demo-go-server"
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	URL             string        `yaml:"url"`               // connection string
	MaxOpenConns    int           `yaml:"max_open_conns"`    // e.g., 10
	MaxIdleConns    int           `yaml:"max_idle_conns"`    // e.g., 5
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"` // e.g., "1h"
}

// CacheConfig holds cache (e.g., Redis) settings.
type CacheConfig struct {
	Host     string `yaml:"host"`     // e.g., "localhost"
	Port     int    `yaml:"port"`     // e.g., 6379
	Password string `yaml:"password"` // if needed
	DB       int    `yaml:"db"`       // e.g., 0
}

// ExternalServicesConfig holds URLs for downstream services.
type ExternalServicesConfig struct {
	UserServiceURL    string `yaml:"user_service_url"`    // e.g., "http://localhost:8080"
	PaymentServiceURL string `yaml:"payment_service_url"` // e.g., "http://localhost:8081"
}

// TimeoutsConfig holds default timeouts for RPCs or handlers.
type TimeoutsConfig struct {
	DefaultRPCTimeout time.Duration `yaml:"default_rpc_timeout"` // e.g., "5s"
	HandlerTimeout    time.Duration `yaml:"handler_timeout"`     // e.g., "10s"
}

// FeatureFlagsConfig holds boolean toggles for feature flags.
type FeatureFlagsConfig struct {
	UseNewAlgorithm bool `yaml:"use_new_algorithm"` // e.g., false in dev
}

// Config is the top-level configuration struct matching dev.yaml/prod.yaml.
type Config struct {
	Environment      string                 `yaml:"environment"` // e.g., "development"
	Server           ServerConfig           `yaml:"server"`
	TLS              TLSConfig              `yaml:"tls"`
	Logging          LoggingConfig          `yaml:"logging"`
	Metrics          MetricsConfig          `yaml:"metrics"`
	Tracing          TracingConfig          `yaml:"tracing"`
	Database         DatabaseConfig         `yaml:"database"`
	Cache            CacheConfig            `yaml:"cache"`
	ExternalServices ExternalServicesConfig `yaml:"external_services"`
	Timeouts         TimeoutsConfig         `yaml:"timeouts"`
	FeatureFlags     FeatureFlagsConfig     `yaml:"feature_flags"`
}

// LoadConfig loads configuration from a YAML file, with optional overrides via environment variables.
// It checks the CONFIG_FILE env var; if unset, uses GO_ENV to pick dev.yaml or prod.yaml under "configs/".
// Returns a Config pointer or an error.
func LoadConfig() (*Config, error) {
	// Determine config file path
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		env := os.Getenv("GO_ENV")
		if env == "production" {
			configPath = "configs/prod.yaml"
		} else {
			// default to development
			configPath = "configs/dev.yaml"
		}
	}
	return loadConfigFromFile(configPath)
}

// loadConfigFromFile reads and parses the YAML at path into Config, then applies environment overrides.
func loadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmtError("failed to read config file %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmtError("failed to parse YAML config %q: %w", path, err)
	}

	// Apply environment variable overrides if set.

	// GRPC_PORT override (overwrites Server.Port)
	if p := os.Getenv("GRPC_PORT"); p != "" {
		if portInt, err := strconv.Atoi(p); err == nil {
			cfg.Server.Port = portInt
		} else {
			return nil, fmtError("invalid GRPC_PORT %q: %w", p, err)
		}
	}
	// GRPC_ENABLE_REFLECTION override (overwrites Server.EnableReflection)
	if refl := os.Getenv("GRPC_ENABLE_REFLECTION"); refl != "" {
		if refl == "true" {
			cfg.Server.EnableReflection = true
		} else if refl == "false" {
			cfg.Server.EnableReflection = false
		} else {
			return nil, fmtError("invalid GRPC_ENABLE_REFLECTION %q: must be \"true\" or \"false\"", refl)
		}
	}
	// LOG_LEVEL override
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		cfg.Logging.Level = lvl
	}
	// METRICS_ENABLED override
	if me := os.Getenv("METRICS_ENABLED"); me != "" {
		if me == "true" {
			cfg.Metrics.Enabled = true
		} else if me == "false" {
			cfg.Metrics.Enabled = false
		} else {
			return nil, fmtError("invalid METRICS_ENABLED %q: must be \"true\" or \"false\"", me)
		}
	}
	// METRICS_PORT override
	if mp := os.Getenv("METRICS_PORT"); mp != "" {
		if portInt, err := strconv.Atoi(mp); err == nil {
			cfg.Metrics.Port = portInt
		} else {
			return nil, fmtError("invalid METRICS_PORT %q: %w", mp, err)
		}
	}
	// DATABASE_URL override
	if dburl := os.Getenv("DATABASE_URL"); dburl != "" {
		cfg.Database.URL = dburl
	}
	// TLS overrides
	if te := os.Getenv("TLS_ENABLED"); te != "" {
		if te == "true" {
			cfg.TLS.Enabled = true
		} else if te == "false" {
			cfg.TLS.Enabled = false
		} else {
			return nil, fmtError("invalid TLS_ENABLED %q: must be \"true\" or \"false\"", te)
		}
	}
	if cert := os.Getenv("TLS_CERT_FILE"); cert != "" {
		cfg.TLS.CertFile = cert
	}
	if key := os.Getenv("TLS_KEY_FILE"); key != "" {
		cfg.TLS.KeyFile = key
	}
	if ca := os.Getenv("TLS_CLIENT_CA_FILE"); ca != "" {
		cfg.TLS.ClientCAFile = ca
	}
	// Additional overrides as needed, e.g., external service URLs:
	if u := os.Getenv("USER_SERVICE_URL"); u != "" {
		cfg.ExternalServices.UserServiceURL = u
	}
	if purl := os.Getenv("PAYMENT_SERVICE_URL"); purl != "" {
		cfg.ExternalServices.PaymentServiceURL = purl
	}
	// Feature flags override example
	if fn := os.Getenv("FEATURE_USE_NEW_ALGORITHM"); fn != "" {
		if fn == "true" {
			cfg.FeatureFlags.UseNewAlgorithm = true
		} else if fn == "false" {
			cfg.FeatureFlags.UseNewAlgorithm = false
		} else {
			return nil, fmtError("invalid FEATURE_USE_NEW_ALGORITHM %q: must be \"true\" or \"false\"", fn)
		}
	}
	// Timeouts overrides (if you wish):
	if dr := os.Getenv("DEFAULT_RPC_TIMEOUT"); dr != "" {
		if d, err := time.ParseDuration(dr); err == nil {
			cfg.Timeouts.DefaultRPCTimeout = d
		} else {
			return nil, fmtError("invalid DEFAULT_RPC_TIMEOUT %q: %w", dr, err)
		}
	}
	if ht := os.Getenv("HANDLER_TIMEOUT"); ht != "" {
		if d, err := time.ParseDuration(ht); err == nil {
			cfg.Timeouts.HandlerTimeout = d
		} else {
			return nil, fmtError("invalid HANDLER_TIMEOUT %q: %w", ht, err)
		}
	}

	return &cfg, nil
}

// fmtError is a helper to format errors.
func fmtError(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
