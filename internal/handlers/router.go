package handlers

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/c4po/tofu-state/internal/config"
	"github.com/c4po/tofu-state/internal/handlers/account"
	"github.com/c4po/tofu-state/internal/handlers/discovery"
	"github.com/c4po/tofu-state/internal/handlers/modules"
	"github.com/c4po/tofu-state/internal/handlers/organizations"
	"github.com/c4po/tofu-state/internal/handlers/state"
	"github.com/c4po/tofu-state/internal/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

func SetupRouter(cfg *config.Config, store sessions.Store, oidcClient *auth.OIDCClient, storageBackend storage.StorageBackend, logger *zap.Logger) *mux.Router {
	logger.Debug("Setting up router")
	router := mux.NewRouter()

	// Service discovery route
	router.HandleFunc("/.well-known/terraform.json", discovery.Handler(logger)).Methods("GET")

	// Auth routes
	router.HandleFunc("/login", auth.HandleLogin(oidcClient, store, logger)).Methods("GET")
	router.HandleFunc("/callback", auth.HandleCallback(oidcClient, store, logger)).Methods("GET")

	// App routes
	appRouter := router.PathPrefix("/app").Subrouter()
	appRouter.Use(auth.AuthMiddleware(store))
	appRouter.HandleFunc("/account", account.DetailsHandler(cfg.BackendURL, logger)).Methods("GET")
	appRouter.HandleFunc("/settings/tokens", auth.HandleTokenRequest(oidcClient, store, logger)).Methods("GET")

	// API routes
	apiRouter := router.PathPrefix("/api/tfe/v2").Subrouter()
	apiRouter.Use(auth.JWTMiddleware(cfg.JWTSecret, logger))

	// Organization routes
	apiRouter.HandleFunc("/organizations/{organization_name}/entitlements",
		organizations.EntitlementsHandler(cfg.BackendURL, logger)).Methods("GET")

	// State management routes
	apiRouter.HandleFunc("/state/{workspace}", state.GetHandler(storageBackend, logger)).Methods("GET")
	apiRouter.HandleFunc("/state/{workspace}", state.PutHandler(storageBackend, logger)).Methods("PUT")

	// Module management routes
	apiRouter.HandleFunc("/modules", modules.ListHandler(storageBackend, logger)).Methods("GET")
	apiRouter.HandleFunc("/modules/{name}", modules.UploadHandler(storageBackend, logger)).Methods("POST")

	// Add catch-all handler
	router.NotFoundHandler = NotFoundHandler(logger)

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
