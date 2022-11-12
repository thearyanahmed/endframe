package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/services/location"
	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"

	"github.com/thearyanahmed/nordsec/pkg/presenter"
)

type nearByRidesUsecase interface {
	// FindNearByRides FindRides @todo use options api instead of map[string]string
	FindNearByRides(ctx context.Context, area locationEntity.Area) ([]locationEntity.Ride, error)
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
	// @todo take from form request
	area := locationEntity.Area{
		X1Y1: locationEntity.Coordinate{
			Lat: 52.3251,
			Lon: 13.453,
		},
		X2Y2: locationEntity.Coordinate{
			Lat: 0,
			Lon: 0,
		},
		X3Y3: locationEntity.Coordinate{
			Lat: 52.3361,
			Lon: 13.475,
		},
		X4Y4: locationEntity.Coordinate{
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
