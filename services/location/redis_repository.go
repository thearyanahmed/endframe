package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
)

// @todo move to dedicated repository directory

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

func (r *RideRepository) getIds(ctx context.Context, ids []string) ([]interface{}, error) {
	return r.datastore.MGet(ctx, ids...).Result()
}

func (r *RideRepository) GetRideEvents(ctx context.Context, geohashKeys []string, keyLen int) ([]RideEventSchema, error) {
	var wg sync.WaitGroup
	wg.Add(keyLen)

	ch := make(chan []string)
	go func() {
		wg.Wait()
		close(ch)
	}()

	for _, hashKey := range geohashKeys {
		go r.getRidesFromGeohash(ctx, ch, &wg, hashKey, 0, -1)
	}

	var collection []RideEventSchema

	for v := range ch {
		events := mapZRangeValueToRideEventSchemaCollection(v)
		collection = append(collection, events...)
	}

	return collection, nil
}

func (r *RideRepository) ApplyStateFilter(ctx context.Context, m map[string]RideEventSchema) map[string]RideEventSchema {
	var keys []string

	for k, _ := range m {
		keys = append(keys, k)
	}

	keyLen := len(keys)
	if keyLen == 0 {
		return m
	}

	res, err := r.getIds(ctx, keys)

	if err != nil {
		return m
	}

	for i := 0; i < keyLen; i++ {
		if res[i] != nil {
			eventSchema := m[keys[i]]
			eventSchema.State = fmt.Sprintf("%v", res[i])
			m[keys[i]] = eventSchema
		}
	}

	return m
}

func (r *RideRepository) GetRideEventsFromMultiGeohash(ctx context.Context, geohashKeys []string) (map[string]RideEventSchema, error) {
	// sanity check
	// if the geohashKeys are empty, then return empty struct
	glen := len(geohashKeys)

	if glen == 0 {
		return map[string]RideEventSchema{}, nil
	}

	// flat array, includes duplicate elements of the same event(ride event schema)
	collection, err := r.GetRideEvents(ctx, geohashKeys, glen)

	if err != nil {
		return map[string]RideEventSchema{}, err
	}

	// from the duplicate ones, filter out
	// @improvement: maybe we can have something like r.applyFilter(function, data)
	m := r.applyUniqueFilter(collection)

	// @todo apply PostFilters // extract to different method
	m = r.ApplyStateFilter(ctx, m)

	return m, nil
}

func (r *RideRepository) applyUniqueFilter(collection []RideEventSchema) map[string]RideEventSchema {
	m := make(map[string]RideEventSchema)

	for _, c := range collection {
		if mv, ok := m[c.RideUuid]; !ok {
			m[c.RideUuid] = c
		} else {
			if mv.Timestamp < c.Timestamp {
				m[c.RideUuid] = c
			}
		}
	}
	return m
}

// @todo find a better name
func (r *RideRepository) getRidesFromGeohash(ctx context.Context, ch chan []string, wg *sync.WaitGroup, geohashKey string, from, till int64) {
	defer wg.Done()

	// @NOTE we are ignoring the error case for this demonstration.
	// Get the data in reverse order
	if result, err := r.datastore.ZRange(ctx, geohashKey, from, till).Result(); err == nil {
		ch <- result
	}
}

// @todo find a better name
func mapZRangeValueToRideEventSchemaCollection(result []string) []RideEventSchema {
	var events []RideEventSchema

	for _, v := range result {
		var ev RideEventSchema

		// @NOTE we are ignoring the error case for this demonstration.
		if err := json.Unmarshal([]byte(v), &ev); err == nil {
			events = append(events, ev)
		}
	}

	return events
}
