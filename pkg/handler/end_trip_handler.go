package handler

import (
	"context"
	"errors"
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
	GetRideEventByUuid(ctx context.Context, rideUuid string) (entity.Event, error)
	TripHasEnded(event entity.Event) bool
	SetRideCurrentStatus(ctx context.Context, event entity.Event) error
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

	rideEvent := eventRequest.ToRideEvent()

	tripEvent, err := h.rideService.GetRideEventByUuid(r.Context(), rideEvent.RideUuid)
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrNotFound(err))
		return
	}

	if tripEvent.TripUuid != eventRequest.TripUuid {
		err = errors.New("given trip uuid does not match current trip uuid")
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if h.rideService.TripHasEnded(tripEvent) {
		presenter.ErrorResponse(w, r, presenter.ErrTripHasAlreadyEnded())
		return
	}

	// @save in database
	recordedEvent, err := h.rideService.RecordEndRideEvent(r.Context(), rideEvent)
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	// add to cool down mode
	err = h.rideService.EnterCooldownMode(r.Context(), recordedEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if err = h.rideService.SetRideCurrentStatus(r.Context(), recordedEvent); err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, presenter.TripEndedResponse(eventRequest.TripUuid))
}
