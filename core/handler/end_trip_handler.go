package handler

import (
	"github.com/thearyanahmed/nordsec/core/presenter"
	"github.com/thearyanahmed/nordsec/core/serializer"
	"github.com/thearyanahmed/nordsec/services/location"
	"net/http"
	"time"
)

type endTripHandler struct {
	locationSvc *location.Service
}

func NewEndTripHandler(locSvc *location.Service) *endTripHandler {
	return &endTripHandler{locationSvc: locSvc}
}

func (h *endTripHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventRequest := &serializer.EndTripRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	// @todo check if trip exists with same status or not
	rideEvent := eventRequest.ToRideEvent()

	// @save in database
	loc, err := h.locationSvc.RecordRideEvent(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	// add to cool down mode
	err = h.locationSvc.StartCooldownForRide(r.Context(), loc.RideUuid, time.Now().Unix(), time.Second*10)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, presenter.TripEndedResponse(eventRequest.TripUuid))
}
