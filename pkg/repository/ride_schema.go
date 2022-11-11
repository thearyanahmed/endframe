package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/thearyanahmed/nordsec/pkg/entity"
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
		RideUuid:  r.UUID,
		Latitude:  r.Coordinates.Latitude,
		Longitude: r.Coordinates.Longitude,
	}
}

func FromRedisGeoLocation(geo redis.GeoLocation) RideLocationSchema {
	coordinate := Coordinates{
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
	}

	return RideLocationSchema{
		UUID:        geo.Name,
		Coordinates: coordinate,
	}
}
