package serializer

import (
	entity2 "github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"github.com/thedevsaddam/govalidator"
)

type StartTripRequest struct {
	RideUuid             string  `json:"ride_uuid" schema:"ride_uuid"`
	ClientUuid           string  `json:"client_uuid" schema:"client_uuid"`
	OriginLatitude       float64 `json:"origin_latitude" schema:"origin_latitude"`
	OriginLongitude      float64 `json:"origin_longitude" schema:"origin_longitude"`
	DestinationLatitude  float64 `json:"destination_latitude" schema:"destination_latitude"`
	DestinationLongitude float64 `json:"destination_longitude" schema:"destination_longitude"`
}

func (r *StartTripRequest) Rules() govalidator.MapData {
	return govalidator.MapData{
		"ride_uuid":             []string{"required", "uuid_v4"},
		"client_uuid":           []string{"required", "uuid_v4"},
		"origin_latitude":       []string{"required", "lat"},
		"origin_longitude":      []string{"required", "lon"},
		"destination_latitude":  []string{"required", "lat"},
		"destination_longitude": []string{"required", "lon"},
	}
}

func (r *StartTripRequest) Origin() entity2.Coordinate {
	return entity2.Coordinate{
		Lat: r.OriginLatitude,
		Lon: r.OriginLongitude,
	}
}

func (r *StartTripRequest) Destination() entity2.Coordinate {
	return entity2.Coordinate{
		Lat: r.DestinationLatitude,
		Lon: r.DestinationLongitude,
	}
}

func (r *StartTripRequest) ToRideEventFromOrigin() entity2.Event {
	return entity2.Event{
		RideUuid:      r.RideUuid,
		Lat:           r.OriginLatitude,
		Lon:           r.OriginLongitude,
		PassengerUuid: r.ClientUuid,
	}
}
