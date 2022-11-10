package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/services/location"
	"net/http"

	"github.com/thearyanahmed/nordsec/core/entity"
	"github.com/thearyanahmed/nordsec/core/presenter"
	"github.com/thearyanahmed/nordsec/core/serializer"
)

type activateRideUsecase interface {
	UpdateRideLocation(ctx context.Context, rideUuid string, lat, long float64) (entity.RideLocationEntity, error)
	FindById(ctx context.Context, uuid string) (entity.RideLocationEntity, error)
}

type activateRideHandler struct {
	usecase     activateRideUsecase
	locationSvc *location.Service
}

func NewActivateRideHandler(usecase activateRideUsecase, locSvc *location.Service) *activateRideHandler {
	return &activateRideHandler{
		usecase:     usecase,
		locationSvc: locSvc,
	}
}

func (h *activateRideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// validate the request
	// check if lat long is valid
	formRequest := &serializer.UpdateRideLocationRequest{}

	if formErrors := serializer.ValidatePostForm(r, formRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	//loc, err := h.usecase.UpdateRideLocation(r.Context(), formRequest.UUID, formRequest.Latitude, formRequest.Longitude)
	//
	//// also need to trigger event
	//if err != nil {
	//	presenter.ErrorResponse(w, r, presenter.FromErr(err))
	//	return
	//}

	cord := location.Coordinate{
		Lat: formRequest.Latitude,
		Lon: formRequest.Longitude,
	}
	loc, err := h.locationSvc.UpdateRideLocation(r.Context(), formRequest.UUID, "", "available", cord)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	//@todo update
	//res := presenter.FromRideLocationEntity(loc)
	presenter.RenderJsonResponse(w, r, http.StatusOK, loc)
}
