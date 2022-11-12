package presenter

import (
	"fmt"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/http"
)

type RideLocationUpdateResponse struct {
	Message string          `json:"message"`
	Event   LocationDetails `json:"event,omitempty"`
}

type LocationDetails struct {
	Uuid      string `json:"uuid"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func FromRideLocationEntity(e entity.Event) RideLocationUpdateResponse {
	return RideLocationUpdateResponse{
		Message: "ride location updated",
		Event: LocationDetails{
			Uuid:      e.RideUuid,
			Latitude:  fmt.Sprintf("%.6f", e.Lat),
			Longitude: fmt.Sprintf("%.6f", e.Lon),
		},
	}
}

func CanNotUpdateLocationViaRiderAppResponse() *Response {
	return &Response{
		HttpStatusCode: http.StatusBadRequest,
		Message:        "can not update ride location via rider app while in route or in cooldown state.",
	}
}
