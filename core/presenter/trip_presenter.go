package presenter

import (
	"fmt"
	"github.com/thearyanahmed/nordsec/services/location"
	"net/http"
)

type TripStarted struct {
	Message string                `json:"message"`
	Route   []location.Coordinate `json:"route"`
	Event   location.RideEvent    `json:"event"`
}

func ErrDistanceTooLowResponse(minDistance float64) *Response {
	return &Response{
		HttpStatusCode: http.StatusBadRequest,
		Message:        fmt.Sprintf("distance too low. minimum distance required %.2f meters or greater", minDistance),
	}
}

func ErrRideUnavailableResponse() *Response {
	return &Response{
		HttpStatusCode: http.StatusUnprocessableEntity,
		Message:        "ride is unavailable",
	}
}

func TripStartedResponse(event location.RideEvent, route []location.Coordinate) TripStarted {
	return TripStarted{
		Message: "trip started",
		Route:   route,
		Event:   event,
	}
}
