package handler

import (
	"net/http"

	"github.com/thearyanahmed/nordsec/core/entity"
	"github.com/thearyanahmed/nordsec/core/presenter"
)

type nearByRidesUsecase interface {
	// @todo use options api instead of map[string]string
	FindRides() ([]entity.RideEntity, error)
}

type nearByRidesHandler struct {
	usecase nearByRidesUsecase
}

func NewNearByRidesHandler(usecase nearByRidesUsecase) *nearByRidesHandler {
	return &nearByRidesHandler{
		usecase: usecase,
	}
}

func (h *nearByRidesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rides, err := h.usecase.FindRides()

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, rides)
}
