package presenter

import (
	"fmt"

	"github.com/thearyanahmed/nordsec/core/entity"
)

type RideLocationUpdateResponse struct {
	Message string          `json:"message"`
	Details LocationDetails `json:"details"`
}

type LocationDetails struct {
	UUID      string `json:"uuid"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func FromRideLocationEntity(e entity.RideLocationEntity) RideLocationUpdateResponse {
	return RideLocationUpdateResponse{
		Message: "ride location updated",
		Details: LocationDetails{
			UUID:      e.UUID,
			Latitude:  fmt.Sprintf("%.8f", e.Latitude),
			Longitude: fmt.Sprintf("%.8f", e.Longitude),
		},
	}
}
