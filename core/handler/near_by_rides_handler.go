package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/services/location"
	"net/http"

	"github.com/thearyanahmed/nordsec/core/entity"
	"github.com/thearyanahmed/nordsec/core/presenter"
)

type nearByRidesUsecase interface {
	// FindRides @todo use options api instead of map[string]string
	FindNearByRides(ctx context.Context, area location.Area) ([]entity.RideEntity, error)
}

type nearByRidesHandler struct {
	usecase nearByRidesUsecase

	// @todo interface
	locationSvc *location.Service
}

func NewNearByRidesHandler(usecase nearByRidesUsecase, loc *location.Service) *nearByRidesHandler {
	return &nearByRidesHandler{
		usecase:     usecase,
		locationSvc: loc,
	}
}

func (h *nearByRidesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	area := location.Area{
		X1Y1: location.Coordinate{
			Lat: 52.3251,
			Lon: 13.453,
		},
		X2Y2: location.Coordinate{
			Lat: 0,
			Lon: 0,
		},
		X3Y3: location.Coordinate{
			Lat: 52.3361,
			Lon: 13.475,
		},
		X4Y4: location.Coordinate{
			Lat: 0,
			Lon: 0,
		},
	}

	//rides, err := h.usecase.FindNearByRides(r.Context(), area)
	rides, err := h.locationSvc.GetRidesInArea(r.Context(), area)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, rides)
}
