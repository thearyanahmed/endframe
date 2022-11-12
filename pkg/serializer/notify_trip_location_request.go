package serializer

import (
	"github.com/thearyanahmed/nordsec/services/location/entity"
	"github.com/thedevsaddam/govalidator"
)

type NotifyTripLocationRequest struct {
	RideUuid      string  `json:"ride_uuid" schema:"ride_uuid"`
	Latitude      float64 `json:"latitude" schema:"latitude"`
	Longitude     float64 `json:"longitude" schema:"longitude"`
	ClientUuid    string  `json:"client_uuid" schema:"client_uuid"`
	PassengerUuid string  `json:"passenger_uuid" schema:"passenger_uuid"`
	TripUuid      string  `json:"trip_uuid" schema:"trip_uuid"`
}

func (r *NotifyTripLocationRequest) Rules() govalidator.MapData {
	return govalidator.MapData{
		"latitude":       []string{"required", "lat"},
		"longitude":      []string{"required", "lon"},
		"trip_uuid":      []string{"required", "uuid_v4"},
		"ride_uuid":      []string{"required", "uuid_v4"},
		"client_uuid":    []string{"required", "uuid_v4"},
		"passenger_uuid": []string{"required", "uuid_v4"},
	}
}

func (r *NotifyTripLocationRequest) ToRideEvent() entity.Event {
	return entity.Event{
		RideUuid:      r.RideUuid,
		Lat:           r.Latitude,
		Lon:           r.Longitude,
		PassengerUuid: r.PassengerUuid,
		TripUuid:      r.TripUuid,
	}
}
