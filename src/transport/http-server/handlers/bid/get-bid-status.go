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

func NewGetBidStatusHandler(logger slog.Logger, uc usecases.GetBidStatusUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "GetBidStatusHandler"
		l := logger.With("op", op)

		bidID := r.PathValue("bidId")
		if bidID == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("missing bidId"))
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("missing username"))
			return
		}

		dto := usecases.GetBidStatusDTO{
			BidID:    bidID,
			Username: username,
		}
		log := l.With("dto", dto)
		status, err := uc.Execute(dto)

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

		api.WriteJSON(w, http.StatusOK, status)
	}
}
