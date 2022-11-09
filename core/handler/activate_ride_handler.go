package handler

import (
	"context"
	"net/http"

	"github.com/thearyanahmed/nordsec/core/presenter"
)

type activateRideUsecase interface {
	UpdateRideLocation(ctx context.Context, rideUuid string, lat, long float64) error
}

type activateRideHandler struct {
	usecase activateRideUsecase
}

func NewActivateRideHandler(usecase activateRideUsecase) *activateRideHandler {
	return &activateRideHandler{
		usecase: usecase,
	}
}

func (h *activateRideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// validate the request
	// check if lat long is valid
	err := h.usecase.UpdateRideLocation(r.Context(), "cities", -76.61219090223312378, 39.29038444452294954)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	res := presenter.RideLocationUpdateResponse{Message: "ride location updated"}
	presenter.RenderJsonResponse(w, r, http.StatusOK, res)
}
