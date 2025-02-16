package main

import (
	"log"
	"net/http"
	"time"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/c4po/tofu-state/internal/config"
	"github.com/c4po/tofu-state/internal/handlers"
	"github.com/c4po/tofu-state/internal/storage"

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
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Initialize OIDC
	oidcClient := auth.NewOIDCClient(cfg.OIDCConfig)

	// Initialize storage
	storageBackend := storage.NewStorageBackend(cfg.StorageConfig)

	// Setup router
	router := handlers.SetupRouter(cfg, store, oidcClient, storageBackend)

	// Add logging middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[DEBUG] %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	log.Printf("Server starting on %s...\n", cfg.ServerAddress)

	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if cfg.CertFile != "" && cfg.KeyFile != "" {
		log.Printf("Using HTTPS with cert: %s and key: %s", cfg.CertFile, cfg.KeyFile)
		log.Fatal(server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile))
	} else {
		log.Println("WARNING: Running in insecure HTTP mode")
		log.Fatal(server.ListenAndServe())
	}
}
