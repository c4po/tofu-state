package state

import (
	"io"
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
	"github.com/c4po/tofu-state/internal/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func GetHandler(backend storage.StorageBackend, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		workspace := vars["workspace"]
		logger.Debug("Getting state", zap.String("workspace", workspace))

		state, err := backend.GetState(workspace)
		if err != nil {
			logger.Error("Failed to retrieve state",
				zap.String("workspace", workspace),
				zap.Error(err))
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to retrieve state",
			})
			return
		}

		logger.Debug("Successfully retrieved state",
			zap.String("workspace", workspace),
			zap.Int("size", len(state)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(state)
	}
}

func PutHandler(backend storage.StorageBackend, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		workspace := vars["workspace"]
		logger.Debug("Putting state", zap.String("workspace", workspace))

		data, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Failed to read request body",
				zap.String("workspace", workspace),
				zap.Error(err))
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusBadRequest,
				Message: "Failed to read request body",
			})
			return
		}

		if err := backend.PutState(workspace, data); err != nil {
			logger.Error("Failed to store state",
				zap.String("workspace", workspace),
				zap.Error(err))
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to store state",
			})
			return
		}

		logger.Debug("Successfully stored state",
			zap.String("workspace", workspace),
			zap.Int("size", len(data)))
		w.WriteHeader(http.StatusOK)
	}
}
