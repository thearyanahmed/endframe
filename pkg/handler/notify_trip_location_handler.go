package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	locationEntity "github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/http"
)

type notifyPositionHandler struct {
	rideService notifyRideService
}

type notifyRideService interface {
	RecordLocationUpdate(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error)
}

func NewNotifyPositionHandler(rideService notifyRideService) *notifyPositionHandler {
	return &notifyPositionHandler{rideService: rideService}
}

func (h *notifyPositionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventRequest := &serializer.NotifyTripLocationRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	// @TODO check in database if trip exists or not.

	rideEvent := eventRequest.ToRideEvent()
	event, err := h.rideService.RecordLocationUpdate(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, event)
}
