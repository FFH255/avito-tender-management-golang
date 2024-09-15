package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"tms/src/core/domain"
	usecases "tms/src/core/services/use-cases/tender"
	"tms/src/pkg/api"
	"tms/src/pkg/logger/sl"
)

func NewChangeTenderStatusHandler(logger slog.Logger, changeTenderStatusUseCase usecases.ChangeTenderStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "ChangeTenderStatusHandler"

		log := logger.With("op", op)

		tenderID := r.PathValue("tenderId")

		if tenderID == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("tenderId is required"))
			log.Error("tenderId is required")
			return
		}

		status := r.URL.Query().Get("status")

		if status == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("status is required"))
			log.Error("status is required")
			return
		}

		username := r.URL.Query().Get("username")

		if username == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("username is required"))
			log.Error("username is required")
			return
		}

		dto := usecases.ChangeTenderStatusDTO{
			TenderID: tenderID,
			Status:   status,
			Username: username,
		}

		log = log.With("dto", dto)

		tender, err := changeTenderStatusUseCase.Execute(dto)

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
			api.WriteJSON(w, http.StatusInternalServerError, api.Error("Internal server error"))
			log.Error("internal server error", sl.Err(err))
			return
		}

		log.Info(fmt.Sprintf("status of tender with id = %s changed", tender.ID))
		api.WriteJSON(w, http.StatusOK, tender)
	}
}
