package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/c4po/tofu-state/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type OIDCClient struct {
	provider *oidc.Provider
	Config   oauth2.Config
}

func NewOIDCClient(cfg config.OIDCConfig) *OIDCClient {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		log.Fatal(err)
	}

	return &OIDCClient{
		provider: provider,
		Config: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "email"},
		},
	}
}

func (c *OIDCClient) VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	verifier := c.provider.Verifier(&oidc.Config{ClientID: c.Config.ClientID})
	return verifier.Verify(ctx, token)
}

func JWTMiddleware(oidc *OIDCClient) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if _, err := oidc.VerifyToken(r.Context(), token); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
