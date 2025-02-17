package modules

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
	"github.com/c4po/tofu-state/internal/storage"
	"go.uber.org/zap"
)

func ListHandler(backend storage.StorageBackend, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Listing modules")

		modules, err := backend.ListModules()
		if err != nil {
			logger.Error("Failed to list modules", zap.Error(err))
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to list modules",
			})
			return
		}

		logger.Debug("Successfully listed modules", zap.Int("count", len(modules)))
		common.RespondWithJSON(w, http.StatusOK, modules)
	}
}

func UploadHandler(backend storage.StorageBackend, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Module upload requested")
		// TODO: Add implementation
	}
}
