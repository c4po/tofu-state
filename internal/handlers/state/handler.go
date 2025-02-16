package state

import (
	"io"
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
	"github.com/c4po/tofu-state/internal/storage"
	"github.com/gorilla/mux"
)

func GetHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		workspace := vars["workspace"]

		state, err := backend.GetState(workspace)
		if err != nil {
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to retrieve state",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(state)
	}
}

func PutHandler(backend storage.StorageBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		workspace := vars["workspace"]

		data, err := io.ReadAll(r.Body)
		if err != nil {
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusBadRequest,
				Message: "Failed to read request body",
			})
			return
		}

		if err := backend.PutState(workspace, data); err != nil {
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to store state",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
