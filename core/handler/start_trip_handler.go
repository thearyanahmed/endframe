package handler

import (
	"fmt"
	"github.com/thearyanahmed/nordsec/core/presenter"
	"github.com/thearyanahmed/nordsec/core/serializer"
	"github.com/thearyanahmed/nordsec/services/location"
	"net/http"
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

func (s *startTripHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	// validation: validate form request, validate the distance.
	// check client is in the same geohash as ride's
	// check if two distance are with in the same range
	tripRequest := &serializer.StartTripRequest{}

	if formErrors := serializer.ValidatePostForm(r, tripRequest); len(formErrors) > 0 {
		presenter.ErrorResponse(w, r, presenter.ErrorValidationFailed(formErrors))
		return
	}

	origin, dest := tripRequest.Origin(), tripRequest.Destination()

	// check the destination first
	if s.locationSvc.DistanceInMeters(origin, dest) <= s.minTripDistance {
		presenter.ErrorResponse(w, r, presenter.ErrDistanceTooLowResponse(s.minTripDistance))
		return
	}

	// check if rideId is in nearby location of the origin
	//ride, err := s.locationSvc.FindRideInLocations(r.Context(), tripRequest.RideUuid, origin)
	//
	//if err != nil {
	//	presenter.ErrorResponse(w, r, presenter.FromErr(err))
	//	return
	//}
	//
	//// @todo refactor
	//if ride.State != "available" {
	//	presenter.ErrorResponse(w, r, presenter.ErrRideUnavailableResponse())
	//	return
	//}

	// @TODO WORK FROM HERE
	// get route
	routes := s.locationSvc.GetRoute(origin, dest, 50)

	fmt.Println("LEN", len(routes))
	presenter.RenderJsonResponse(w, r, http.StatusOK, routes)
	// ride goes to origin
	// start notify event with current passenger id
	//

	// response: {
	//	message: 'trip started',
	//  route: [{lat1,lon1}, {lat2,lon2}...{latN,lonN}]
}
