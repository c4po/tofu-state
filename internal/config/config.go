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
	CertFile      string
	KeyFile       string
	OIDCConfig    OIDCConfig
	StorageConfig storage.StorageConfig
	ExternalHost  string
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
		port = "8080" // Default port
	}

	// Validate port number
	if portNum, _ := strconv.ParseUint(port, 10, 16); portNum < 1 || portNum > 65535 {
		log.Fatalf("Invalid port number: %s (must be 1-65535)", port)
	}

	serverAddress := net.JoinHostPort("0.0.0.0", port)

	if clientID == "" || clientSecret == "" {
		log.Fatal("OIDC_CLIENT_ID and OIDC_CLIENT_SECRET must be set")
	}

	cfg := &Config{
		ServerAddress: serverAddress,
		CertFile:      certFile,
		KeyFile:       keyFile,
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

	// Validate redirect port is within allowed range
	redirectPort := 10000 // First port in range
	if _, port, err := net.SplitHostPort(cfg.OIDCConfig.RedirectURL); err == nil {
		if p, err := strconv.Atoi(port); err == nil {
			redirectPort = p
		}
	}
	if redirectPort < 10000 || redirectPort > 10010 {
		log.Fatal("OIDC redirect port must be between 10000-10010")
	}

	return cfg
}
