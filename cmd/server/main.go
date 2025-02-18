package main

import (
	"net/http"
	"time"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/c4po/tofu-state/internal/config"
	"github.com/c4po/tofu-state/internal/handlers"
	"github.com/c4po/tofu-state/internal/storage"
	"github.com/spf13/viper"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	config.LoadConfig()

	if len(viper.GetString("session_secret")) < 32 {
		logger.Fatal("Session secret must be at least 32 characters long")
	}

	store := sessions.NewCookieStore([]byte(viper.GetString("session_secret")))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Initialize OIDC
	oidcClient := auth.NewOIDCClient(logger)

	// Initialize storage
	storageBackend := storage.NewStorageBackend(logger)

	// Setup router
	router := handlers.SetupRouter(store, oidcClient, storageBackend, logger)

	// Add logging middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("Request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
			next.ServeHTTP(w, r)
		})
	})

	logger.Info("Server starting on", zap.String("address", viper.GetString("server_address")))

	server := &http.Server{
		Addr:         viper.GetString("server_address"),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if viper.GetString("cert_file") != "" && viper.GetString("key_file") != "" {
		logger.Info("Using HTTPS with cert", zap.String("cert", viper.GetString("cert_file")), zap.String("key", viper.GetString("key_file")))
		err := server.ListenAndServeTLS(viper.GetString("cert_file"), viper.GetString("key_file"))
		if err != nil {
			logger.Fatal("Failed to start HTTPS server", zap.Error(err))
		}
	} else {
		logger.Warn("WARNING: Running in insecure HTTP mode")
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}
}
