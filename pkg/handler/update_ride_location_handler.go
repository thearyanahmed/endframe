package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/services/location"
	entity2 "github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"
	"time"

	"github.com/thearyanahmed/nordsec/pkg/entity"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
)

type activateRideUsecase interface {
	UpdateRideLocation(ctx context.Context, rideUuid string, lat, long float64) (entity.RideLocationEntity, error)
	FindById(ctx context.Context, uuid string) (entity.RideLocationEntity, error)
}

type updateRideLocationHandler struct {
	usecase activateRideUsecase
	// @todo extract interface
	locationSvc *location.Service
}

func NewUpdateRideLocationHandler(usecase activateRideUsecase, locSvc *location.Service) *updateRideLocationHandler {
	return &updateRideLocationHandler{
		usecase:     usecase,
		locationSvc: locSvc,
	}
}

func (h *updateRideLocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// validate the request
	// check if lat long is valid
	// @todo find a better name
	eventRequest := &serializer.RecordRideEventRequest{}

	if formErrors := serializer.ValidatePostForm(r, eventRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	rideEvent := entity2.RideEvent{
		RideUuid:      eventRequest.RideUuid,
		Lat:           eventRequest.Latitude,
		Lon:           eventRequest.Longitude,
		PassengerUuid: "",
		State:         "available", // @todo handle this
		Timestamp:     time.Now().Unix(),
	}

	loc, err := h.locationSvc.RecordRideEvent(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	//@todo update
	//res := presenter.FromRideLocationEntity(loc)
	presenter.RenderJsonResponse(w, r, http.StatusOK, loc)
}
