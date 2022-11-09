package handler

import (
	"context"
	"net/http"

	"github.com/thearyanahmed/nordsec/core/entity"
	"github.com/thearyanahmed/nordsec/core/presenter"
)

type activateRideUsecase interface {
	UpdateRideLocation(ctx context.Context, rideUuid string, lat, long float64) (entity.RideLocationEntity, error)
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
	loc, err := h.usecase.UpdateRideLocation(r.Context(), "uuid-1234-5678", -76.61219090223312378, 39.29038444452294954)

	// also need to trigger event
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	res := presenter.FromRideLocationEntity(loc)
	presenter.RenderJsonResponse(w, r, http.StatusOK, res)
}
