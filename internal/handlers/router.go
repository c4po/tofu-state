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
)

func SetupRouter(cfg *config.Config, store sessions.Store, oidcClient *auth.OIDCClient, storageBackend storage.StorageBackend) *mux.Router {
	router := mux.NewRouter()

	// Service discovery route
	router.HandleFunc("/.well-known/terraform.json", discovery.Handler(cfg)).Methods("GET")

	// Auth routes
	router.HandleFunc("/login", auth.HandleLogin(oidcClient, store)).Methods("GET")
	router.HandleFunc("/callback", auth.HandleCallback(oidcClient, store)).Methods("GET")

	// App routes
	appRouter := router.PathPrefix("/app").Subrouter()
	appRouter.Use(auth.AuthMiddleware(store))
	appRouter.HandleFunc("/account", account.DetailsHandler(cfg.BackendURL)).Methods("GET")
	appRouter.HandleFunc("/settings/tokens", auth.HandleTokenRequest(oidcClient, store)).Methods("GET")

	// API routes
	apiRouter := router.PathPrefix("/api/tfe/v2").Subrouter()
	apiRouter.Use(auth.JWTMiddleware(cfg.JWTSecret))

	// Organization routes
	apiRouter.HandleFunc("/organizations/{organization_name}/entitlements",
		organizations.EntitlementsHandler(cfg.BackendURL)).Methods("GET")

	// State management routes
	apiRouter.HandleFunc("/state/{workspace}", state.GetHandler(storageBackend)).Methods("GET")
	apiRouter.HandleFunc("/state/{workspace}", state.PutHandler(storageBackend)).Methods("PUT")

	// Module management routes
	apiRouter.HandleFunc("/modules", modules.ListHandler(storageBackend)).Methods("GET")
	apiRouter.HandleFunc("/modules/{name}", modules.UploadHandler(storageBackend)).Methods("POST")

	// Add catch-all handler
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	return router
}
