package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/http"
)

type endTripHandler struct {
	rideService endTripRideService
}

type endTripRideService interface {
	RecordEndRideEvent(ctx context.Context, event entity.Event) (entity.Event, error)
	EnterCooldownMode(ctx context.Context, event entity.Event) error
}

func NewEndTripHandler(riderSvc endTripRideService) *endTripHandler {
	return &endTripHandler{rideService: riderSvc}
}

func (h *endTripHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventRequest := &serializer.EndTripRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	// @TODO check if trip exists with same status or not
	rideEvent := eventRequest.ToRideEvent()

	// @save in database

	recordedEvent, err := h.rideService.RecordEndRideEvent(r.Context(), rideEvent)
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	// add to cool down mode
	err = h.rideService.EnterCooldownMode(r.Context(), recordedEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, presenter.TripEndedResponse(eventRequest.TripUuid))
}
