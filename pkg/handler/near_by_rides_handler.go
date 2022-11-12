package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
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
	filterRequest := &serializer.NearByRidesRequest{}

	if formErrors := serializer.ValidateGetQuery(r, filterRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	rides, err := h.usecase.FindNearByRides(r.Context(), filterRequest.ToArea(r))

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, rides)
}
