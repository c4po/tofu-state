package handlers

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/storage"
)

func GetStateHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement state retrieval
	}
}

func PutStateHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement state storage
	}
}
