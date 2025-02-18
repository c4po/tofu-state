package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func forwardHandler(logger *zap.Logger) http.HandlerFunc {
	tfeBackendURL := viper.GetString("tfe_backend_url")
	target, err := url.Parse(tfeBackendURL)
	if err != nil {
		logger.Fatal("Invalid backend URL", zap.Error(err))
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(w http.ResponseWriter, r *http.Request) {
		// Remove the API prefix from the path before forwarding
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/tfe/v2")

		// Update the Host header to match the target
		r.Host = target.Host

		logger.Debug("Forwarding request",
			zap.String("path", r.URL.Path),
			zap.String("target", target.String()))

		proxy.ServeHTTP(w, r)
	}
}
