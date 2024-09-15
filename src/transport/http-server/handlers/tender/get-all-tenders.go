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

func NewGetAllTendersHandler(log slog.Logger, getAllTendersUseCase usecases.GetAllTendersUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "GetAllTendersHandler"
		l := log.With("op", op)

		limit, _ := api.ParseIntQueryParam(r, "limit")
		offset, _ := api.ParseIntQueryParam(r, "offset")
		serviceType := api.ParseStringQueryParam(r, "service_type")

		dto := usecases.GetAllTendersDTO{
			Limit:       limit,
			Offset:      offset,
			ServiceType: serviceType,
		}

		l = l.With("dto", dto)

		tenders, err := getAllTendersUseCase.Execute(dto)

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
			l.Error("error while executing GetAllTendersUseCase", "error", err)
			resp := api.Error("Внутренняя ошибка сервера")
			api.WriteJSON(w, http.StatusInternalServerError, resp)
			return
		}

		l.Info("GetAllTendersUseCase executed successfully", tenders)

		api.WriteJSON(w, http.StatusOK, tenders)
	}
}
