package handlers

import (
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"tms/src/core/domain"
	usecases "tms/src/core/services/use-cases/tender"
	"tms/src/pkg/api"
	"tms/src/pkg/logger/sl"
)

type EditTenderHandlerBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ServiceType *string `json:"serviceType"`
}

func NewEditTenderHandler(logger slog.Logger, editTenderUseCase usecases.EditTenderUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "EditTenderHandler"

		log := logger.With("op", op)

		tenderID := r.PathValue("tenderId")

		if tenderID == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("tenderId is required"))
			log.Error("tenderId is required")
			return
		}

		username := r.URL.Query().Get("username")

		if username == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("username is required"))
			log.Error("username is required")
			return
		}

		body, err := api.ReadJSON[EditTenderHandlerBody](r)

		if err != nil {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("cannot parse body"))
			log.Error("cannot parse body", sl.Err(err))
			return
		}

		dto := usecases.EditTenderUseCaseDTO{
			TenderID:    tenderID,
			Username:    username,
			Name:        body.Name,
			Description: body.Description,
			ServiceType: body.ServiceType,
		}

		log = log.With("dto", dto)

		tender, err := editTenderUseCase.Execute(dto)

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
			log.Error("cannot execute editTenderUseCase", sl.Err(err))
			return
		}

		api.WriteJSON(w, http.StatusOK, tender)
	}
}
