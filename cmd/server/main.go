package main

import (
	"log"
	"net/http"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/c4po/tofu-state/internal/config"
	"github.com/c4po/tofu-state/internal/handlers"
	"github.com/c4po/tofu-state/internal/storage"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	cfg := config.LoadConfig()

	if len(cfg.SessionSecret) < 32 {
		log.Fatal("Session secret must be at least 32 characters long")
	}

	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   cfg.Environment == "production",
		SameSite: http.SameSiteLaxMode,
	}

	router := mux.NewRouter()

	// Add logging middleware to log all requests
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[DEBUG] %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Initialize OIDC
	oidcClient := auth.NewOIDCClient(cfg.OIDCConfig)

	// Initialize storage
	storageBackend := storage.NewStorageBackend(cfg.StorageConfig)

	// Service discovery route
	router.HandleFunc("/.well-known/terraform.json", handlers.DiscoveryHandler(cfg)).Methods("GET")

	// Auth routes
	router.HandleFunc("/login", handlers.HandleLogin(oidcClient, store)).Methods("GET")
	router.HandleFunc("/callback", handlers.HandleCallback(oidcClient, store)).Methods("GET")

	appRouter := router.PathPrefix("/app").Subrouter()
	appRouter.Use(auth.AuthMiddleware(store))
	appRouter.HandleFunc("/settings/tokens", handlers.HandleTokenRequest(oidcClient, store)).Methods("GET")

	// Protected routes
	apiRouter := router.PathPrefix("/api/tfe/v2").Subrouter()
	apiRouter.Use(auth.JWTMiddleware(cfg.JWTSecret))
	apiRouter.HandleFunc("/account/details", handlers.AccountDetailsHandler()).Methods("GET")

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
