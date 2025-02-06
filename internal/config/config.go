package config

import (
	"log"
	"os"

	"github.com/c4po/tofu-state/internal/storage"
)

type Config struct {
	ServerAddress string
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

	if clientID == "" || clientSecret == "" {
		log.Fatal("OIDC_CLIENT_ID and OIDC_CLIENT_SECRET must be set")
	}

	return &Config{
		ServerAddress: ":8080",
		CertFile:      certFile,
		KeyFile:       keyFile,
		OIDCConfig: OIDCConfig{
			IssuerURL:    "https://accounts.google.com",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  "http://localhost:8080/callback",
		},
		StorageConfig: storage.StorageConfig{
			Type:      "local",
			LocalPath: "./data",
		},
	}
}
