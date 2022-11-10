package serializer

import (
	"github.com/thedevsaddam/govalidator"
)

type UpdateRideLocationRequest struct {
	UUID      string  `json:"uuid" schema:"uuid"`
	Latitude  float64 `json:"latitude" schema:"latitude"`
	Longitude float64 `json:"longitude" schema:"longitude"`
}

func (r *UpdateRideLocationRequest) Rules() govalidator.MapData {
	return govalidator.MapData{
		"uuid":      []string{"required", "uuid_v4"},
		"latitude":  []string{"required", "lat"},
		"longitude": []string{"required", "lon"},
	}
}
