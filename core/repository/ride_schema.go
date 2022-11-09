package repository

import (
	"github.com/thearyanahmed/nordsec/core/entity"
)

type RideLocationSchema struct {
	UUID        string      `bson:"_id"`
	Coordinates Coordinates `bson:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
}

func (r *RideLocationSchema) ToEntity() entity.RideLocationEntity {
	return entity.RideLocationEntity{
		UUID:      r.UUID,
		Latitude:  r.Coordinates.Latitude,
		Longitude: r.Coordinates.Longitude,
	}
}
