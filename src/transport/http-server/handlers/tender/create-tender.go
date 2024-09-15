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

func NewCreateTenderHandler(log slog.Logger, createTenderUseCase usecases.CreateTenderUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "CreateTenderHandler"

		l := log.With("op", op)

		body, err := api.ReadJSON[usecases.CreateTenderDTO](r)

		if err != nil {
			api.WriteJSON(w, http.StatusBadRequest, api.Error(fmt.Sprintf("unable to parse request: %s", err.Error())))
			l.Error("unable to parse request", err.Error())
			return
		}

		l = l.With("body", body)

		tender, err := createTenderUseCase.Execute(*body)

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
			api.WriteJSON(w, http.StatusBadRequest, api.Error(fmt.Sprintf("unable to create Tender: %s", err.Error())))
			l.Error("unable to create Tender", err.Error())
			return
		}

		l.Info("tender created", tender)
		api.WriteJSON(w, http.StatusOK, tender)
	}
}
