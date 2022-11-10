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

func (r *RideRepository) FindById(ctx context.Context, uuid string) (RideLocationSchema, error) {
	//var rideLocation RideLocationSchema

	//if err := r.datastore.FindOne(ctx, bson.M{"_id": uuid}).Decode(&rideLocation); err != nil {
	return RideLocationSchema{}, nil
	//}

	//return rideLocation, nil
}

func (r *RideRepository) UpdateLocation(ctx context.Context, uuid string, lat, lon float64) (RideLocationSchema, error) {
	loc := redis.GeoLocation{
		Name:      uuid,
		Latitude:  lat,
		Longitude: lon,
	}

	// @todo need to update list name
	_, err := r.datastore.GeoAdd(ctx, "test", &loc).Result()

	if err != nil {
		return RideLocationSchema{}, err
	}

	return FromRedisGeoLocation(loc), nil
}
