package handler

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	entity2 "github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/http"
)

type startTripHandler struct {
	rideService rideService
}

type rideService interface {
	GetMinimumTripDistance() float64
	IsRideAvailable(ride entity2.Ride) bool
	GetRoute(origin, dest entity2.Coordinate) []entity2.Coordinate
	RecordNewRideEvent(ctx context.Context, event entity2.Event) (entity2.Event, error)
	DistanceIsGreaterThanMinimumDistance(origin, destination entity2.Coordinate) bool
	FindRideInLocation(ctx context.Context, rideUuid string, rideLocation entity2.Coordinate) (entity2.Ride, error)
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

	if !h.rideService.IsRideAvailable(ride) {
		presenter.ErrorResponse(w, r, presenter.ErrRideUnavailableResponse())
		return
	}

	routes := h.rideService.GetRoute(origin, dest)
	rideEvent := tripRequest.ToRideEventFromOrigin()

	event, err := h.rideService.RecordNewRideEvent(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusCreated, presenter.TripStartedResponse(event, routes))
}
