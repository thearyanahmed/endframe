package serializer

import (
	"github.com/thearyanahmed/nordsec/services/location/entity"
	"github.com/thedevsaddam/govalidator"
)

type RecordRideEventRequest struct {
	RideUuid  string  `json:"ride_uuid" schema:"ride_uuid"`
	Latitude  float64 `json:"latitude" schema:"latitude"`
	Longitude float64 `json:"longitude" schema:"longitude"`
}

func (r *RecordRideEventRequest) Rules() govalidator.MapData {
	return govalidator.MapData{
		"ride_uuid": []string{"required", "uuid_v4"},
		"latitude":  []string{"required", "lat"},
		"longitude": []string{"required", "lon"},
	}
}

func (r *RecordRideEventRequest) ToLocationCoordinate() entity.Coordinate {
	return entity.Coordinate{
		Lat: r.Latitude,
		Lon: r.Longitude,
	}
}
