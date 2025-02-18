package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	log.Println("Starting to load configuration...")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // path to look for the config file in
	// Set environment variable prefix
	viper.SetEnvPrefix("TFS")
	viper.AutomaticEnv() // read in environment variables that match

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			log.Fatalf("Error reading config file: %s", err)
		}
		// Config file not found; ignore error if desired
		log.Printf("No config file found. Using defaults and environment variables")
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set default values
	viper.SetDefault("jwt_secret", "default-insecure-secret-please-change-in-prod")
	viper.SetDefault("session_secret", "default-insecure-secret-please-change-in-prod")
	viper.SetDefault("tfe_backend_url", "http://localhost:8080")

	// Set default server address based on SSL configuration
	if viper.GetString("server_address") == "" {
		if viper.GetString("cert_file") != "" && viper.GetString("key_file") != "" {
			viper.Set("server_address", "0.0.0.0:443")
		} else {
			viper.Set("server_address", "0.0.0.0:80")
		}
	}

	log.Printf("Configuration loaded successfully:")
	log.Printf("- Server Address: %s", viper.GetString("server_address"))
	log.Printf("- TFE Backend URL: %s", viper.GetString("tfe_backend_url"))
	log.Printf("- SSL Enabled: %v", viper.GetString("cert_file") != "" && viper.GetString("key_file") != "")
	log.Printf("- OIDC Issuer URL: %s", viper.GetString("oidc.issuer_url"))
}
