package handlers

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/c4po/tofu-state/internal/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

func SetupRouter(store sessions.Store, oidcClient *auth.OIDCClient, storageBackend storage.StorageBackend, logger *zap.Logger) *mux.Router {
	logger.Debug("Setting up router")
	router := mux.NewRouter()

	apiPath := "/api/tfe/v2"

	// Service discovery route
	router.HandleFunc("/.well-known/terraform.json", DiscoveryHandler(apiPath, logger)).Methods("GET")

	// Auth routes
	router.HandleFunc("/login", auth.HandleLogin(oidcClient, store, logger)).Methods("GET")
	router.HandleFunc("/callback", auth.HandleCallback(oidcClient, store, logger)).Methods("GET")

	// App routes
	appRouter := router.PathPrefix("/app").Subrouter()
	appRouter.Use(auth.AuthMiddleware(store))
	appRouter.HandleFunc("/settings/tokens", auth.HandleTokenRequest(oidcClient, store, logger)).Methods("GET")
	appRouter.HandleFunc("/account", forwardHandler(logger)).Methods("GET")

	// API routes
	apiRouter := router.PathPrefix(apiPath).Subrouter()
	apiRouter.HandleFunc("/ping", PingHandler(logger)).Methods("GET")

	// Single forward handler for all API routes
	apiRouter.PathPrefix("/").Handler(forwardHandler(logger))

	// Add logging middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr))
			next.ServeHTTP(w, r)
		})
	})

	logger.Info("Router setup complete")
	return router
}
