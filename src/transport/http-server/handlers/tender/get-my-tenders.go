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

func NewGetMyTendersHandlers(log slog.Logger, getUserTendersUseCase usecases.GetUserTendersUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "GetMyTendersHandler"

		l := log.With("op", op)

		limit, _ := api.ParseIntQueryParam(r, "limit")
		offset, _ := api.ParseIntQueryParam(r, "offset")
		username := api.ParseStringQueryParam(r, "username")

		if username == nil {
			l.Error("username is required")
			api.WriteJSON(w, http.StatusBadRequest, api.Error("username is required"))
			return
		}

		dto := usecases.GetUserTendersDTO{
			Limit:    limit,
			Offset:   offset,
			Username: *username,
		}

		l = l.With("dto", dto)

		tenders, err := getUserTendersUseCase.Execute(dto)

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
			l.Error("failed to execute getUserTendersUseCase", err.Error())
			api.WriteJSON(w, http.StatusInternalServerError, api.Error("internal server error"))
			return
		}

		api.WriteJSON(w, http.StatusOK, tenders)
	}
}
