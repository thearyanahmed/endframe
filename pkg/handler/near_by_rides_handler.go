package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/http"

	"github.com/thearyanahmed/nordsec/pkg/presenter"
)

type nearByRidesUsecase interface {
	FindNearByRides(ctx context.Context, area entity.Area) ([]entity.Ride, error)
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
	// @todo take from form request
	area := entity.Area{
		X1Y1: entity.Coordinate{
			Lat: 52.3251,
			Lon: 13.453,
		},
		X2Y2: entity.Coordinate{
			Lat: 0,
			Lon: 0,
		},
		X3Y3: entity.Coordinate{
			Lat: 52.3361,
			Lon: 13.475,
		},
		X4Y4: entity.Coordinate{
			Lat: 0,
			Lon: 0,
		},
	}

	rides, err := h.usecase.FindNearByRides(r.Context(), area)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, rides)
}
