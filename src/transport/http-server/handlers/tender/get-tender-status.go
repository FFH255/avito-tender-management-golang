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

func NewGetTenderStatus(log slog.Logger, getTenderStatusUseCase usecases.GetTenderStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "getTenderStatus"

		l := log.With("op", op)

		tenderID := r.PathValue("tenderId")

		username := api.ParseStringQueryParam(r, "username")

		if tenderID == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("missing parameter: tenderId"))
			l.Error("missing parameter: tenderId")
			return
		}

		if username == nil {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("missing parameter: username"))
			l.Error("missing parameter: username")
			return
		}

		dto := usecases.GetTenderStatusDTO{
			TenderID: tenderID,
			Username: *username,
		}

		l = l.With("dto", dto)

		status, err := getTenderStatusUseCase.Execute(dto)

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
			l.Error("internal server error", sl.Err(err))
			return
		}

		l.Info("get Tender status success")
		api.WriteJSON(w, http.StatusOK, status)
	}
}
