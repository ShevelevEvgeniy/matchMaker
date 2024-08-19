package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	DTOs "matchMaker/internal/dto"
	"matchMaker/internal/http_server/api/v1/handlers/response"
)

type Service interface {
	SaveUsers(ctx context.Context, user DTOs.User) error
}

type MatchMakerHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   Service
}

func NewMatchMakerHandler(log *zap.Logger, validator *validator.Validate, service Service) *MatchMakerHandler {
	return &MatchMakerHandler{
		log:       log,
		validator: validator,
		service:   service,
	}
}

func (h *MatchMakerHandler) SaveUsers(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.log.Info("Received HTTP POST request: " + r.RequestURI)

		var dto DTOs.User
		err := json.NewDecoder(r.Body).Decode(&dto)
		if err != nil {
			h.log.Error("failed to decode user", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.BadRequest("failed to decode user"))
			return
		}

		err = h.validator.Struct(dto)
		if err != nil {
			h.log.Error("failed to validate user", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.BadRequest(err.Error()))
			return
		}

		err = h.service.SaveUsers(ctx, dto)
		if err != nil {
			h.log.Error("failed to add user", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.InternalServerError())
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Created())
	}
}
