package auth

import (
	"context"
	"log"

	"github.com/c4po/tofu-state/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
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
