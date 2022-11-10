package location

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mmcloughlin/geohash"
)

type Service struct {
	geohashLength uint
	repo          rideRepository
}

type rideRepository interface {
	UpdateLocation(ctx context.Context, ghash string, trip RideEventSchema) error
	GetRideEventsFromMultiGeohash(ctx context.Context, geohashKeys []string) (interface{}, error)
}

// NewLocationService
// @todo use a builder pattern to build out the service?
func NewLocationService(ds *redis.Client) *Service {
	repo := NewRideRepository(ds, "trips_test_01")

	return &Service{
		geohashLength: uint(6),
		repo:          repo,
	}
}

func getGeoHash(lat, lon float64, precision uint) string {
	return geohash.EncodeWithPrecision(lat, lon, precision)
}

func (s *Service) RecordRideEvent(ctx context.Context, event RideEvent) (RideEvent, error) {
	// get the location geohash
	ghash := getGeoHash(event.Lat, event.Lon, s.geohashLength)

	err := geohash.Validate(ghash)

	if err != nil {
		return RideEvent{}, err
	}

	loc := fromRideEventEntity(event).WithNewUuid()

	fmt.Println("EVENT", event, "Ghash", ghash)
	err = s.repo.UpdateLocation(ctx, ghash, *loc)

	if err != nil {
		return RideEvent{}, err
	}

	return loc.ToEntity(), nil
}

func (s *Service) GetRidesInArea(ctx context.Context, area Area) (interface{}, error) {
	// first get the data in area
	// then apply the filters if any
	neighbours := area.GetNeighbourGeohashFromCenter(s.geohashLength)

	fmt.Println("AREA", area)
	fmt.Println("neighbours", neighbours)

	rides, err := s.repo.GetRideEventsFromMultiGeohash(ctx, neighbours)

	if err != nil {
		fmt.Println("ERROR ->", err)
		return map[string]RideEventSchema{}, err
	}

	return rides, nil
	//return fromRideEventCollection(rides), nil
}
