package modules

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
	"github.com/c4po/tofu-state/internal/storage"
)

func ListHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modules, err := backend.ListModules()
		if err != nil {
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to list modules",
			})
			return
		}
		common.RespondWithJSON(w, http.StatusOK, modules)
	}
}

func UploadHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement module upload
		// TODO: Add implementation
	}
}
