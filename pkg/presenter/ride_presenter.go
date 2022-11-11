package presenter

import (
	"fmt"
	"github.com/thearyanahmed/nordsec/services/location/entity"
)

type RideLocationUpdateResponse struct {
	Message string          `json:"message"`
	Event   LocationDetails `json:"event"`
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
