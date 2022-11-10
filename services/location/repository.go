package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
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

	_, err = r.datastore.ZAdd(ctx, ghash, &redis.Z{
		Score:  float64(trip.Timestamp),
		Member: jsonStr,
	}).Result()

	return err
}

func (r *RideRepository) getIds(ctx context.Context, ids []string) {
	data := r.datastore.MGet(ctx, ids...).Val()
	fmt.Println("MGET() DATA", data)
}

func (r *RideRepository) GetRideEventsFromMultiGeohash(ctx context.Context, geohashKeys []string) ([]RideEventSchema, error) {
	// get all rides from redis
	// use concurrency with waitGroups

	// @todo remove these
	geohashKeys = []string{}
	geohashKeys = append(geohashKeys, "u339gf")

	glen := len(geohashKeys)

	if glen == 0 {
		return []RideEventSchema{}, nil
	}

	ch := make(chan []string)

	var wg sync.WaitGroup

	wg.Add(glen)
	go func() {
		wg.Wait()
		close(ch)
	}()

	//u339gf
	for _, hashKey := range geohashKeys {
		go r.getRidesFromGeohash(ctx, ch, &wg, hashKey, 0, -1)
	}

	var events []RideEventSchema

	for v := range ch {
		event := mapZRangeValueToRideEventSchemaCollection(v)
		events = append(events, event...)
	}

	return events, nil
}

func getZrange() {
	//result, err := r.datastore.ZRange(ctx, "u339gf", 0, -1).Result()
	//
	//if err != nil {
	//	return []RideEventSchema{}, err
	//}
	//
	//fmt.Println("RESULT->", result)
	//
	//events := mapZRangeValueToRideEventSchemaCollection(result)
}

func (r *RideRepository) getRidesFromGeohash(ctx context.Context, ch chan []string, wg *sync.WaitGroup, geohashKey string, from, till int64) {
	defer wg.Done()

	result, err := r.datastore.ZRange(ctx, geohashKey, from, till).Result()

	if err != nil {
		return
	}

	fmt.Println("RESULT->", result)
	ch <- result
}

func mapZRangeValueToRideEventSchemaCollection(result []string) []RideEventSchema {
	var events []RideEventSchema

	for _, v := range result {
		var x RideEventSchema
		err := json.Unmarshal([]byte(v), &x)

		if err != nil {
			continue
		}

		events = append(events, x)
	}

	return events
}
