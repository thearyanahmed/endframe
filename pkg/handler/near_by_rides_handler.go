package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/services/location"
	entity2 "github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"

	"github.com/thearyanahmed/nordsec/pkg/entity"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
)

type nearByRidesUsecase interface {
	// FindRides @todo use options api instead of map[string]string
	FindNearByRides(ctx context.Context, area entity2.Area) ([]entity.RideEntity, error)
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
	area := entity2.Area{
		X1Y1: entity2.Coordinate{
			Lat: 52.3251,
			Lon: 13.453,
		},
		X2Y2: entity2.Coordinate{
			Lat: 0,
			Lon: 0,
		},
		X3Y3: entity2.Coordinate{
			Lat: 52.3361,
			Lon: 13.475,
		},
		X4Y4: entity2.Coordinate{
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
