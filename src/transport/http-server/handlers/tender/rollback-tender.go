package handlers

import (
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"strconv"
	"tms/src/core/domain"
	usecases "tms/src/core/services/use-cases/tender"
	"tms/src/pkg/api"
	"tms/src/pkg/logger/sl"
)

func NewRollbackTenderHandler(logger slog.Logger, rollbackTenderUseCase usecases.RollBackTenderUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "RollbackTenderHandler"

		log := logger.With("op", op)

		tenderID := r.PathValue("tenderId")

		if tenderID == "" {
			api.WriteJSON(w, http.StatusBadRequest, "tenderId is required")
			log.Error("tenderId is required")
			return
		}

		versionStr := r.PathValue("version")

		if versionStr == "" {
			api.WriteJSON(w, http.StatusBadRequest, "version is required")
			log.Error("version is required")
			return
		}

		version, err := strconv.Atoi(versionStr)

		if err != nil {
			api.WriteJSON(w, http.StatusBadRequest, "version should be int")
			log.Error("version should be int")
			return
		}

		username := r.URL.Query().Get("username")

		if username == "" {
			api.WriteJSON(w, http.StatusBadRequest, "username is required")
			log.Error("username is required")
			return
		}

		dto := usecases.RollBackTenderUseCaseDTO{
			TenderID: tenderID,
			Version:  version,
			Username: username,
		}

		log = log.With("dto", dto)

		tender, err := rollbackTenderUseCase.Execute(dto)

		if err != nil {
			if errors.Is(errors.Cause(err), domain.ErrValidation) {
				api.WriteJSON(w, http.StatusBadRequest, api.Error(err.Error()))
				log.Error("validation failed", sl.Err(err))
				return
			}
			if errors.Is(errors.Cause(err), domain.ErrNotFound) {
				api.WriteJSON(w, http.StatusBadRequest, api.Error(err.Error()))
				log.Error("some entity not found", sl.Err(err))
				return
			}
			if errors.Is(errors.Cause(err), domain.ErrAlreadyExist) {
				api.WriteJSON(w, http.StatusBadRequest, api.Error(err.Error()))
				log.Error("already exists", sl.Err(err))
				return
			}
			if errors.Is(errors.Cause(err), domain.ErrNoPermission) {
				api.WriteJSON(w, http.StatusForbidden, api.Error(err.Error()))
				log.Error("permission denied", sl.Err(err))
				return
			}
			if errors.Is(errors.Cause(err), domain.ErrUserNotFound) {
				api.WriteJSON(w, http.StatusUnauthorized, api.Error(err.Error()))
				log.Error("user not found", sl.Err(err))
				return
			}
			api.WriteJSON(w, http.StatusInternalServerError, api.Error("internal server error"))
			log.Error("failed to rollback tender", sl.Err(err))
			return
		}

		api.WriteJSON(w, http.StatusOK, tender)
		log.Info("rolled back tender")
	}
}
