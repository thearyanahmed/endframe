package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/http"
)

type nearByRidesUsecase interface {
	FindNearByRides(ctx context.Context, area entity.Area, stateFilter string) ([]entity.Ride, error)
}

type nearByRidesHandler struct {
	riderService nearByRidesUsecase
}

func NewNearByRidesHandler(rideSvc nearByRidesUsecase) *nearByRidesHandler {
	return &nearByRidesHandler{
		riderService: rideSvc,
	}
}

func (h *nearByRidesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filterRequest := &serializer.NearByRidesRequest{}

	if formErrors := serializer.ValidateGetQuery(r, filterRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	rides, err := h.riderService.FindNearByRides(r.Context(), filterRequest.ToArea(r), r.URL.Query().Get("state"))

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, rides)
}
