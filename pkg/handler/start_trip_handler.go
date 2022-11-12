package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"
)

type startTripHandler struct {
	rideService rideService
}

type rideService interface {
	GetMinimumTripDistance() float64
	RideIsAvailable(ride entity.Ride) bool
	GetRoute(origin, dest entity.Coordinate) []entity.Coordinate
	RecordRideEvent(ctx context.Context, event entity.Event) (entity.Event, error)
	DistanceIsGreaterThanMinimumDistance(origin, destination entity.Coordinate) bool
	FindRideInLocation(ctx context.Context, rideUuid string, rideLocation entity.Coordinate) (entity.Ride, error)
}

func NewStartTripHandler(rideService rideService) *startTripHandler {
	return &startTripHandler{
		rideService: rideService,
	}
}

func (h *startTripHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tripRequest := &serializer.StartTripRequest{}

	if formErrors := serializer.ValidatePostForm(r, tripRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	origin, dest := tripRequest.Origin(), tripRequest.Destination()

	if !h.rideService.DistanceIsGreaterThanMinimumDistance(origin, dest) {
		presenter.ErrorResponse(w, r, presenter.ErrDistanceTooLowResponse(h.rideService.GetMinimumTripDistance()))
		return
	}

	// @TODO | check if rideId is in nearby location of the origin
	ride, err := h.rideService.FindRideInLocation(r.Context(), tripRequest.RideUuid, origin)
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	if !h.rideService.RideIsAvailable(ride) {
		presenter.ErrorResponse(w, r, presenter.ErrRideUnavailableResponse())
		return
	}

	routes := h.rideService.GetRoute(origin, dest)
	rideEvent := tripRequest.ToRideEventFromOrigin()

	event, err := h.rideService.RecordRideEvent(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	// start notify event with current passenger id
	presenter.RenderJsonResponse(w, r, http.StatusCreated, presenter.TripStartedResponse(event, routes))
}
