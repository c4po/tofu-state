package backend

import (
	"context"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/viper"
)

// GetTFEClient returns a configured TFE client
func GetTFEClient(token string) (*tfe.Client, error) {
	config := &tfe.Config{
		Address: viper.GetString("tfe_backend_url"),
		Token:   token,
		// Enable retrying on server errors
		RetryServerErrors: true,
	}

	return tfe.NewClient(config)
}

// ValidateClient checks if the client can connect to the TFE server
func ValidateClient(client *tfe.Client) error {
	// Try to read the current user as a basic validation
	_, err := client.Users.ReadCurrent(context.Background())
	return err
}
