package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"
)

type activateRideUsecase interface {
	// Updated
	UpdateRideLocation(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error)

	// @TODO Old
	//FindById(ctx context.Context, uuid string) (entity.RideLocationEntity, error)
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

	// @todo | TASK get ride, check if it's in route or in cooldown
	rideEvent := eventRequest.ToRideEvent().SetStateAsRoaming().SetCurrentTimestamp()

	loc, err := h.usecase.UpdateRideLocation(r.Context(), *rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	res := presenter.FromRideLocationEntity(loc)
	presenter.RenderJsonResponse(w, r, http.StatusOK, res)
}
