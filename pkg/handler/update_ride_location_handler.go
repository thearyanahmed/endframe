package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/http"
)

type activateRideUsecase interface {
	UpdateRideLocation(ctx context.Context, event entity.Event) (entity.Event, error)
	CanBeUpdatedViaRiderApp(ctx context.Context, rideUuid string, loc entity.Coordinate) (bool, error)
	SetRideCurrentStatus(ctx context.Context, event entity.Event) error
	GetRideEventByUuid(ctx context.Context, rideUuid string) (entity.Event, error)
	IsInRoute(state string) bool
}

type updateRideLocationHandler struct {
	rideService activateRideUsecase
}

func NewUpdateRideLocationHandler(usecase activateRideUsecase) *updateRideLocationHandler {
	return &updateRideLocationHandler{
		rideService: usecase,
	}
}

func (h *updateRideLocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventRequest := &serializer.RecordRideEventRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	// @NOTE:
	// As of now, there is no persistent storage.
	// In a real life scenario there will be a persistent storage where we can make a query and see the ride's
	// current stations.
	if tripEvent, err := h.rideService.GetRideEventByUuid(r.Context(), eventRequest.RideUuid); err == nil && h.rideService.IsInRoute(tripEvent.State) {
		presenter.ErrorResponse(w, r, presenter.ErrNotFound())
		return
	}

	rideEvent := eventRequest.ToRideEvent().SetStateAsRoaming().SetCurrentTimestamp()

	recordedEvent, err := h.rideService.UpdateRideLocation(r.Context(), *rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if err = h.rideService.SetRideCurrentStatus(r.Context(), recordedEvent); err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, presenter.FromRideLocationEntity(recordedEvent))
}
