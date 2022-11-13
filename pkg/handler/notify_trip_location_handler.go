package handler

import (
	"context"
	"errors"
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
	SetRideCurrentStatus(ctx context.Context, event locationEntity.Event) error
	GetRideEventByUuid(ctx context.Context, rideUuid string) (locationEntity.Event, error)
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

	tripEvent, err := h.rideService.GetRideEventByUuid(r.Context(), eventRequest.RideUuid)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrNotFound())
		return
	}

	if tripEvent.TripUuid != eventRequest.TripUuid {
		err = errors.New("given trip uuid does not match current trip uuid")
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	rideEvent := eventRequest.ToRideEvent()
	event, err := h.rideService.RecordLocationUpdate(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if err = h.rideService.SetRideCurrentStatus(r.Context(), event); err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, event)
}
