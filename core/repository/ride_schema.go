package repository

import "github.com/go-redis/redis/v8"

type RideSchema struct {
	UUID      string
	Latitude  float64
	Longitude float64
}

func FromRedisGeoLocation(geo redis.GeoLocation) RideSchema {
	return RideSchema{
		UUID:      geo.Name,
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
	}
}
