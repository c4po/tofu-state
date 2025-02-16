package config

import (
	"log"
	"os"

	"github.com/c4po/tofu-state/internal/storage"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerAddress string                `yaml:"server_address,omitempty"`
	SessionSecret string                `yaml:"session_secret,omitempty"`
	JWTSecret     string                `yaml:"jwt_secret,omitempty"`
	CertFile      string                `yaml:"cert_file,omitempty"`
	KeyFile       string                `yaml:"key_file,omitempty"`
	OIDCConfig    OIDCConfig            `yaml:"oidc"`
	StorageConfig storage.StorageConfig `yaml:"storage"`
	BackendURL    string                `yaml:"backend_url,omitempty"`
}

type OIDCConfig struct {
	IssuerURL    string `yaml:"issuer_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

func LoadConfig() *Config {
	// Load base config from YAML
	cfg := &Config{}
	if data, err := os.ReadFile("config.yaml"); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			log.Fatal("Error parsing config.yaml: ", err)
		}
	}

	// Override with environment variables (TFS_ prefix)
	err := envconfig.Process("TFS", cfg)
	if err != nil {
		log.Fatal("Error processing environment variables: ", err)
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "default-insecure-secret-please-change-in-prod"
	}

	if cfg.SessionSecret == "" {
		cfg.SessionSecret = "default-insecure-secret-please-change-in-prod"
	}

	if cfg.ServerAddress == "" {
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			cfg.ServerAddress = "0.0.0.0:443"
		} else {
			cfg.ServerAddress = "0.0.0.0:80"
		}
	}

	if cfg.BackendURL == "" {
		cfg.BackendURL = "http://localhost:8080"
	}

	return cfg
}
