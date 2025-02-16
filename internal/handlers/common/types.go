package common

import (
	"context"
	"net/http"
	"time"

	"github.com/c4po/tofu-state/internal/backend"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/jsonapi"
)

const (
	DefaultTimeout = 10 * time.Second
)

type ContextKey string

const (
	UserEmailKey ContextKey = "userEmail"
	UserTokenKey ContextKey = "userToken"
)

type HandlerError struct {
	Status  int
	Message string
}

func (e HandlerError) Error() string {
	return e.Message
}

// Common handler utilities
func GetTFEClientFromContext(r *http.Request, backendURL string) (*tfe.Client, error) {
	token, ok := r.Context().Value(UserTokenKey).(string)
	if !ok || token == "" {
		return nil, HandlerError{
			Status:  http.StatusUnauthorized,
			Message: "No authentication token found",
		}
	}

	client, err := backend.GetTFEClient(backendURL, token)
	if err != nil {
		return nil, HandlerError{
			Status:  http.StatusInternalServerError,
			Message: "Failed to initialize TFE client",
		}
	}

	return client, nil
}

func WithTimeout(parent context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, duration)
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(status)
	return jsonapi.MarshalPayload(w, payload)
}

func RespondWithError(w http.ResponseWriter, err error) {
	if handlerErr, ok := err.(HandlerError); ok {
		http.Error(w, handlerErr.Message, handlerErr.Status)
		return
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
