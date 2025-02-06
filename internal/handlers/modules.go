package handlers

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/storage"
)

func ListModulesHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement module listing
	}
}

func UploadModuleHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement module upload
	}
}
