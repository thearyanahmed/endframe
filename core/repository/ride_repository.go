package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RideRepository struct {
	datastore *redis.Client
}

func NewRideRepository(ds *redis.Client) *RideRepository {
	return &RideRepository{
		datastore: ds,
	}
}

func (r *RideRepository) UpdateLocation(ctx context.Context, uuid string, lat, long float64) (RideLocationSchema, error) {
	loc := redis.GeoLocation{
		Name:      uuid,
		Latitude:  lat,
		Longitude: long,
	}

	// @todo need to update list name
	_, err := r.datastore.GeoAdd(ctx, "test", &loc).Result()

	if err != nil {
		return RideLocationSchema{}, err
	}

	return FromRedisGeoLocation(loc), err
}
