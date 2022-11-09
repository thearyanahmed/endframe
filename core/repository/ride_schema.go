package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/thearyanahmed/nordsec/core/entity"
)

type RideLocationSchema struct {
	UUID      string
	Latitude  float64
	Longitude float64
}

func FromRedisGeoLocation(geo redis.GeoLocation) RideLocationSchema {
	return RideLocationSchema{
		UUID:      geo.Name,
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
	}
}

func (r *RideLocationSchema) ToEntity() entity.RideLocationEntity {
	return entity.RideLocationEntity{
		UUID:      r.UUID,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}
