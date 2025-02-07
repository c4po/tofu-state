package config

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/c4po/tofu-state/internal/storage"
)

type Config struct {
	ServerAddress string
	Environment   string
	SessionSecret string
	JWTSecret     string
	CertFile      string
	KeyFile       string
	OIDCConfig    OIDCConfig
	StorageConfig storage.StorageConfig
}

type OIDCConfig struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func LoadConfig() *Config {
	clientID := os.Getenv("OIDC_CLIENT_ID")
	clientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	port := os.Getenv("PORT")

	if port == "" {
		if certFile != "" && keyFile != "" {
			port = "443" // Default HTTPS port
		} else {
			port = "80" // Default HTTP port
		}
	}

	// Validate port number
	if portNum, _ := strconv.ParseUint(port, 10, 16); portNum < 1 || portNum > 65535 {
		log.Fatalf("Invalid port number: %s (must be 1-65535)", port)
	}

	serverAddress := net.JoinHostPort("0.0.0.0", port)

	if clientID == "" || clientSecret == "" {
		log.Fatal("OIDC_CLIENT_ID and OIDC_CLIENT_SECRET must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-insecure-secret-please-change-in-prod"
	}

	cfg := &Config{
		ServerAddress: serverAddress,
		CertFile:      certFile,
		KeyFile:       keyFile,
		JWTSecret:     jwtSecret,
		SessionSecret: "abcdefghijklmnopqrstuvwxyz123456",
		OIDCConfig: OIDCConfig{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  "/callback",
		},
		StorageConfig: storage.StorageConfig{
			Type:      "local",
			LocalPath: "./data",
		},
	}

	return cfg
}
