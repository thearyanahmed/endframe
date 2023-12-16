package handler

import (
	"context"
	"net/http"

	"github.com/thearyanahmed/endframe/pkg/presenter"
	"github.com/thearyanahmed/endframe/pkg/serializer"
	"github.com/thearyanahmed/endframe/pkg/service/location/entity"
)

type startTripHandler struct {
	rideService rideService
}

type rideService interface {
	GetMinimumTripDistance() float64
	IsRideAvailable(ride entity.Ride) bool
	GetRoute(origin, dest entity.Coordinate) []entity.Coordinate
	RecordNewRideEvent(ctx context.Context, event entity.Event) (entity.Event, error)
	DistanceIsGreaterThanMinimumDistance(origin, destination entity.Coordinate) bool
	FindRideInLocation(ctx context.Context, rideUuid string, rideLocation entity.Coordinate) (entity.Ride, error)
	SetRideCurrentStatus(ctx context.Context, event entity.Event) error
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

	ride, err := h.rideService.FindRideInLocation(r.Context(), tripRequest.RideUuid, origin)
	if err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
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
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	if err = h.rideService.SetRideCurrentStatus(r.Context(), event); err != nil {
		presenter.ErrorResponse(w, r, presenter.ErrFrom(err))
		return
	}

	presenter.RenderJsonResponse(w, r, http.StatusCreated, presenter.TripStartedResponse(event, routes))
}
