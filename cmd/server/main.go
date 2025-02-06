package main

import (
	"log"
	"net/http"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/c4po/tofu-state/internal/config"
	"github.com/c4po/tofu-state/internal/handlers"
	"github.com/c4po/tofu-state/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	router := mux.NewRouter()

	// Initialize OIDC
	oidcClient := auth.NewOIDCClient(cfg.OIDCConfig)

	// Initialize storage
	storageBackend := storage.NewStorageBackend(cfg.StorageConfig)

	// Service discovery route
	router.HandleFunc("/.well-known/terraform.json", handlers.DiscoveryHandler(cfg)).Methods("GET")

	// Auth routes
	router.HandleFunc("/login", handlers.HandleLogin(oidcClient)).Methods("GET")
	router.HandleFunc("/callback", handlers.HandleCallback(oidcClient)).Methods("GET")

	// Protected routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(auth.JWTMiddleware(oidcClient))

	// State management
	apiRouter.HandleFunc("/state/{workspace}", handlers.GetStateHandler(storageBackend)).Methods("GET")
	apiRouter.HandleFunc("/state/{workspace}", handlers.PutStateHandler(storageBackend)).Methods("PUT")

	// Module management
	apiRouter.HandleFunc("/modules", handlers.ListModulesHandler(storageBackend)).Methods("GET")
	apiRouter.HandleFunc("/modules/{name}", handlers.UploadModuleHandler(storageBackend)).Methods("POST")

	log.Printf("Server starting on %s...\n", cfg.ServerAddress)

	if cfg.CertFile != "" && cfg.KeyFile != "" {
		log.Printf("Using HTTPS with cert: %s and key: %s", cfg.CertFile, cfg.KeyFile)
		log.Fatal(http.ListenAndServeTLS(
			cfg.ServerAddress,
			cfg.CertFile,
			cfg.KeyFile,
			router,
		))
	} else {
		log.Println("WARNING: Running in insecure HTTP mode")
		log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
	}
}
