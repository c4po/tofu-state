package auth

import (
	"context"

	"github.com/c4po/tofu-state/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type OIDCClient struct {
	provider *oidc.Provider
	Config   oauth2.Config
	logger   *zap.Logger
}

func NewOIDCClient(cfg config.OIDCConfig, logger *zap.Logger) *OIDCClient {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		logger.Fatal("Failed to create OIDC provider", zap.Error(err))
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
		logger: logger,
	}
}

func (c *OIDCClient) VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	verifier := c.provider.Verifier(&oidc.Config{ClientID: c.Config.ClientID})
	c.logger.Debug("Verifying OIDC token", zap.String("client_id", c.Config.ClientID))
	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		c.logger.Error("Failed to verify OIDC token", zap.Error(err))
		return nil, err
	}
	c.logger.Debug("Successfully verified OIDC token")
	return idToken, nil
}
