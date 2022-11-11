package handler

import (
	"github.com/google/uuid"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/services/location"
	"github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"
	"time"
)

type startTripHandler struct {
	// @todo extract interface
	minTripDistance float64 // minimum distance between origin and destination, in meters
	locationSvc     *location.Service
}

func NewStartTripHandler(service *location.Service) *startTripHandler {
	return &startTripHandler{
		locationSvc:     service,
		minTripDistance: 500, // in meters
	}
}

func (h *startTripHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	// validation: validate form request, validate the distance.
	// check client is in the same geohash as ride'h
	// check if two distance are with in the same range
	tripRequest := &serializer.StartTripRequest{}

	if formErrors := serializer.ValidatePostForm(r, tripRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	origin, dest := tripRequest.Origin(), tripRequest.Destination()

	// check the destination first
	if h.locationSvc.DistanceInMeters(origin, dest) <= h.minTripDistance {
		presenter.ErrorResponse(w, r, presenter.ErrDistanceTooLowResponse(h.minTripDistance))
		return
	}

	// @TODO enable the following
	// check if rideId is in nearby location of the origin
	ride, err := h.locationSvc.FindRideInLocations(r.Context(), tripRequest.RideUuid, origin)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	// @todo refactor
	if ride.State != "available" {
		presenter.ErrorResponse(w, r, presenter.ErrRideUnavailableResponse())
		return
	}

	// @TODO WORK FROM HERE
	// get route
	routes := h.locationSvc.GetRoute(origin, dest, 50)

	// @Todo add to database
	rideEvent := entity.RideEvent{
		RideUuid:      tripRequest.RideUuid,
		Lat:           origin.Lat,
		Lon:           origin.Lon,
		PassengerUuid: tripRequest.ClientUuid,
		State:         "in_route", // @todo handle this
		Timestamp:     time.Now().Unix(),
		TripUuid:      uuid.New().String(),
	}

	event, err := h.locationSvc.RecordRideEvent(r.Context(), rideEvent)

	if err != nil {
		presenter.ErrorResponse(w, r, presenter.FromErr(err))
		return
	}

	// start notify event with current passenger id
	presenter.RenderJsonResponse(w, r, http.StatusCreated, presenter.TripStartedResponse(event, routes))
}
