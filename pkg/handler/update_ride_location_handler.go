package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"
)

type activateRideUsecase interface {
	UpdateRideLocation(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error)
	CanBeUpdatedViaRiderApp(ctx context.Context, rideUuid string, loc locationEntity.Coordinate) (bool, error)
}

type updateRideLocationHandler struct {
	usecase activateRideUsecase
}

func NewUpdateRideLocationHandler(usecase activateRideUsecase) *updateRideLocationHandler {
	return &updateRideLocationHandler{
		usecase: usecase,
	}
}

func (h *updateRideLocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventRequest := &serializer.RecordRideEventRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	// @NOTE:
	// The following section has been commented out. Because, as of now, there is no persistent storage.
	// In a real life scenario there will be a persistent storage where we can make a query and see the ride's
	// current stations.

	// @TODO implement update ride's status in REDIS ride:status:uuid roaming, in-route

	//can, err := h.usecase.CanBeUpdatedViaRiderApp(r.Context(), eventRequest.RideUuid, eventRequest.ToLocationCoordinate())
	//
	//if err != nil {
	//	presenter.ErrorResponse(w, r, presenter.FromErr(err))
	//	return
	//}
	//
	//if !can {
	//	presenter.ErrorResponse(w, r, presenter.CanNotUpdateLocationViaRiderAppResponse())
	//	return
	//}

	rideEvent := eventRequest.ToRideEvent().SetStateAsRoaming().SetCurrentTimestamp()

	loc, err := h.usecase.UpdateRideLocation(r.Context(), *rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusOK, presenter.FromRideLocationEntity(loc))
}
