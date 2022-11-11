package presenter

import (
	"fmt"

	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"
)

type RideLocationUpdateResponse struct {
	Message string          `json:"message"`
	Event   LocationDetails `json:"event"`
}

type LocationDetails struct {
	UUID      string `json:"uuid"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func FromRideLocationEntity(e locationEntity.Event) RideLocationUpdateResponse {
	return RideLocationUpdateResponse{
		Message: "ride location updated",
		Event: LocationDetails{
			UUID:      e.RideUuid,
			Latitude:  fmt.Sprintf("%.5f", e.Lat),
			Longitude: fmt.Sprintf("%.5f", e.Lon),
		},
	}
}
