package repository

import "github.com/go-redis/redis/v8"

type RideRepository struct {
	datastore *redis.Client
}

func NewRideRepository(ds *redis.Client) *RideRepository {
	return &RideRepository{
		datastore: ds,
	}
}
