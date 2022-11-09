package presenter

import (
	"github.com/thearyanahmed/nordsec/core/entity"
)

type RideLocationUpdateResponse struct {
	Message string `json:"message"`
	Details LocationDetails
}

type LocationDetails struct {
	UUID      string  `json:"uuid"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func FromRideLocationEntity(e entity.RideLocationEntity) RideLocationUpdateResponse {
	return RideLocationUpdateResponse{
		Message: "ride location updated",
		Details: LocationDetails{
			UUID:      e.UUID,
			Latitude:  e.Latitude,
			Longitude: e.Longitude,
		},
	}
}
