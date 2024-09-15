package handlers

import (
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"tms/src/core/domain"
	usecases "tms/src/core/services/use-cases/bid"
	"tms/src/pkg/api"
	"tms/src/pkg/logger/sl"
)

func NewCreateBidHandler(logger slog.Logger, createBidUseCase usecases.CreateBidUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "CreateBidHandler"
		log := logger.With("op", op)

		body, err := api.ReadJSON[usecases.CreateBidDTO](r)
		if err != nil {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("can not read body"))
			log.Error("can not read body", "err", err)
			return
		}
		log = log.With("body", body)

		bid, err := createBidUseCase.Execute(*body)
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

		api.WriteJSON(w, http.StatusOK, bid)
		log.Info("Create bid success")
	}
}
