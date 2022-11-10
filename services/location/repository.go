package location

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RideRepository struct {
	datastore *redis.Client
	ridesKey  string
}

func NewRideRepository(ds *redis.Client, ridesKey string) *RideRepository {
	return &RideRepository{
		datastore: ds,
		ridesKey:  ridesKey,
	}
}

func (r *RideRepository) UpdateLocation(ctx context.Context, ghash string, trip RideEventSchema) error {
	jsonStr, err := json.Marshal(trip)

	if err != nil {
		return err
	}
	fmt.Println(ghash)

	_, err = r.datastore.ZAdd(ctx, ghash, &redis.Z{
		Score:  float64(trip.Timestamp),
		Member: jsonStr,
	}).Result()

	return err
}

func (r *RideRepository) getIds(ctx context.Context, ids []string) {
	data := r.datastore.MGet(ctx, ids...).Val()
	fmt.Println("MGET DATA", data)
}

func (r *RideRepository) GetRideEventsFromMultiGeohash(ctx context.Context, geohashKeys []string) ([]RideEventSchema, error) {
	return []RideEventSchema{}, errors.New("yet to be implemented")
}
