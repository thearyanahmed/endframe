package location

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mmcloughlin/geohash"
	"github.com/umahmood/haversine"
	"time"
)

type Service struct {
	geohashLength uint
	repo          rideRepository
}

type rideRepository interface {
	UpdateLocation(ctx context.Context, ghash string, trip RideEventSchema) error
	GetRideEventsFromMultiGeohash(ctx context.Context, geohashKeys []string) (map[string]RideEventSchema, error)
	ApplyStateFilter(ctx context.Context, m map[string]RideEventSchema) map[string]RideEventSchema
	SetToCooldown(ctx context.Context, details RideCooldownEvent) error
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

	rides, err := s.repo.GetRideEventsFromMultiGeohash(ctx, neighbours)

	if err != nil {
		return map[string]RideEventSchema{}, err
	}

	rides = s.repo.ApplyStateFilter(ctx, rides)

	return rides, nil
	//return fromRideEventCollection(rides), nil
}

func (s *Service) DistanceInMeters(a, b Coordinate) float64 {
	origin := haversine.Coord{Lat: a.Lat, Lon: a.Lon}
	dest := haversine.Coord{Lat: b.Lat, Lon: b.Lon}

	_, km := haversine.Distance(origin, dest)

	return km * 1000 // convert to meters
}

// @change ride event schema to entity
func (s *Service) FindRideInLocations(ctx context.Context, rideUuid string, origin Coordinate) (Ride, error) {
	// get origin
	// get neighbours
	ghash := geohash.EncodeWithPrecision(origin.Lat, origin.Lon, s.geohashLength)

	// get all rides in neighbours
	neighbours := geohash.Neighbors(ghash)

	// check if rides there
	rides, err := s.repo.GetRideEventsFromMultiGeohash(ctx, neighbours)
	if err != nil {
		return Ride{}, err
	}

	rideEvent, ok := rides[rideUuid]

	// if it doesn't exist, return error
	if !ok {
		return Ride{}, errors.New("ride not in nearby location")
	}

	// but if it exists, apply the state filter, so if any vehicle is in cool down mode, we can validate it correctly
	m := map[string]RideEventSchema{rideUuid: rideEvent}

	rides = s.repo.ApplyStateFilter(ctx, m)

	rideEventSchema := rides[rideUuid]

	return rideEventSchema.ToRideEntity(), nil
}

func (s *Service) GetRoute(origin, destination Coordinate, intervalPoints int) []Coordinate {
	x1, x2 := origin.Lat, destination.Lat
	y1, y2 := origin.Lon, destination.Lon

	dx := (x2 - x1) / float64(intervalPoints)
	dy := (y2 - y1) / float64(intervalPoints)

	var route []Coordinate

	route = append(route, origin)

	for i := 1; i < intervalPoints-1; i++ {
		point := Coordinate{
			Lat: x1 + float64(i)*dx,
			Lon: y1 + float64(i)*dy,
		}
		route = append(route, point)
	}

	route = append(route, destination)

	return route
}

func (s *Service) StartCooldownForRide(ctx context.Context, rideUuid string, timestamp int64, duration time.Duration) error {
	ev := RideCooldownEvent{
		RideUuid:  rideUuid,
		Timestamp: timestamp,
		Duration:  duration,
	}

	return s.repo.SetToCooldown(ctx, ev)
}
