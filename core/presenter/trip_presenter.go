package presenter

import (
	"fmt"
	"github.com/thearyanahmed/nordsec/services/location/entity"
	"net/http"
)

type TripStarted struct {
	Message string              `json:"message"`
	Route   []entity.Coordinate `json:"route"`
	Event   entity.RideEvent    `json:"event"`
}

type TripEnded struct {
	Message  string `json:"message"`
	TripUuid string `json:"trip_uuid"`
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

func TripStartedResponse(event entity.RideEvent, route []entity.Coordinate) TripStarted {
	return TripStarted{
		Message: "trip started",
		Route:   route,
		Event:   event,
	}
}

func TripEndedResponse(tripUuid string) TripEnded {
	return TripEnded{
		Message:  "trip has ended",
		TripUuid: tripUuid,
	}
}
