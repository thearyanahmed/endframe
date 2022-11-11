package handler

import (
	"github.com/thearyanahmed/nordsec/core/presenter"
	"github.com/thearyanahmed/nordsec/core/serializer"
	"github.com/thearyanahmed/nordsec/services/location"
	"net/http"
)

type notifyPositionHandler struct {
	locationSvc *location.Service
}

func NewNotifyPositionHandler(locSvc *location.Service) *notifyPositionHandler {
	return &notifyPositionHandler{locationSvc: locSvc}
}

func (h *notifyPositionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// validate request
	// add event with in_route
	eventRequest := &serializer.NotifyTripLocationRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	// @TODO check in database if trip exists or not.
	rideEvent := eventRequest.ToRideEvent()

	loc, err := h.locationSvc.RecordRideEvent(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	// return response
	presenter.RenderJsonResponse(w, r, http.StatusOK, loc)
}
