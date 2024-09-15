package handlers

import (
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"strconv"
	"tms/src/core/domain"
	usecases "tms/src/core/services/use-cases/bid"
	"tms/src/pkg/api"
	"tms/src/pkg/logger/sl"
)

func NewRollBackHandler(logger slog.Logger, uc usecases.RollbackBidUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bidID := r.PathValue("bidId")
		if bidID == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("bidID is required"))
			logger.Error("bidID is required")
			return
		}

		versionStr := r.PathValue("version")
		if versionStr == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("version is required"))
			logger.Error("version is required")
			return
		}

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("version should be a number"))
			logger.Error("version should be a number")
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			api.WriteJSON(w, http.StatusBadRequest, api.Error("username is required"))
			logger.Error("username is required")
			return
		}

		dto := usecases.RollbackBidDTO{
			BidID:    bidID,
			Version:  version,
			Username: username,
		}
		log := logger.With("dto", dto)
		bid, err := uc.Execute(dto)
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
	}
}
