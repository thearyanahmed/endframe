package location

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/mmcloughlin/geohash"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"github.com/thearyanahmed/nordsec/pkg/service/location/repository"
	"github.com/thearyanahmed/nordsec/pkg/service/location/schema"
	"github.com/umahmood/haversine"
	"time"
)

var RideNotInLocation = errors.New("ride not in nearby location")

type Service struct {
	geohashLength uint
	repo          rideRepository
}

type rideRepository interface {
	UpdateLocation(ctx context.Context, ghash string, trip schema.RideEventSchema) error
	GetRideEventsFromMultiGeohash(ctx context.Context, geohashKeys []string) (map[string]schema.RideEventSchema, error)
	ApplyCooldownStateFilter(ctx context.Context, m map[string]schema.RideEventSchema) map[string]schema.RideEventSchema
	SetToCooldown(ctx context.Context, details schema.RideCooldownEvent) error
}

func NewLocationService(ds *redis.Client, redisKey string) *Service {
	repo := repository.NewRideRepository(ds, redisKey)
	return &Service{
		geohashLength: uint(6), // could have taken from config as well
		repo:          repo,
	}
}

func getGeoHash(lat, lon float64, precision uint) string {
	return geohash.EncodeWithPrecision(lat, lon, precision)
}

func (s *Service) RecordRideEvent(ctx context.Context, event entity.Event) (entity.Event, error) {
	ghash := getGeoHash(event.Lat, event.Lon, s.geohashLength)

	err := geohash.Validate(ghash)

	if err != nil {
		return entity.Event{}, err
	}

	loc := schema.FromRideEventEntity(event).WithNewUuid()

	err = s.repo.UpdateLocation(ctx, ghash, *loc)

	if err != nil {
		return entity.Event{}, err
	}

	return loc.ToEntity(), nil
}

func (s *Service) GetRidesInArea(ctx context.Context, area entity.Area) ([]entity.Ride, error) {
	neighbours := area.GetNeighbourGeohashFromCenter(s.geohashLength)

	rides, err := s.repo.GetRideEventsFromMultiGeohash(ctx, neighbours)

	if err != nil {
		return []entity.Ride{}, err
	}

	rides = s.repo.ApplyCooldownStateFilter(ctx, rides)

	return schema.FromRideEventCollectionToEntity(rides), nil
}

func (s *Service) FindRideInLocation(ctx context.Context, rideUuid string, origin entity.Coordinate) (entity.Ride, error) {
	ghash := geohash.EncodeWithPrecision(origin.Lat, origin.Lon, s.geohashLength)

	neighbours := geohash.Neighbors(ghash)

	rides, err := s.repo.GetRideEventsFromMultiGeohash(ctx, neighbours)
	if err != nil {
		return entity.Ride{}, err
	}

	rideEvent, ok := rides[rideUuid]

	if !ok {
		return entity.Ride{}, RideNotInLocation
	}

	m := map[string]schema.RideEventSchema{rideUuid: rideEvent}

	rides = s.repo.ApplyCooldownStateFilter(ctx, m)

	rideEventSchema := rides[rideUuid]

	return rideEventSchema.ToRideEntity(), nil
}

func (s *Service) StartCooldownForRide(ctx context.Context, rideUuid string, timestamp int64, duration time.Duration) error {
	ev := schema.RideCooldownEvent{
		RideUuid:  rideUuid,
		Timestamp: timestamp,
		Duration:  duration,
	}

	return s.repo.SetToCooldown(ctx, ev)
}

// GetRoute returns an array of points, step by step coordinates from origin to destination.
func (s *Service) GetRoute(origin, destination entity.Coordinate, intervalPoints int) []entity.Coordinate {
	x1, x2 := origin.Lat, destination.Lat
	y1, y2 := origin.Lon, destination.Lon

	dx := (x2 - x1) / float64(intervalPoints)
	dy := (y2 - y1) / float64(intervalPoints)

	var route []entity.Coordinate

	route = append(route, origin)

	for i := 1; i < intervalPoints-1; i++ {
		point := entity.Coordinate{
			Lat: x1 + float64(i)*dx,
			Lon: y1 + float64(i)*dy,
		}
		route = append(route, point)
	}

	route = append(route, destination)

	return route
}

func (s *Service) DistanceInMeters(a, b entity.Coordinate) float64 {
	origin := haversine.Coord{Lat: a.Lat, Lon: a.Lon}
	dest := haversine.Coord{Lat: b.Lat, Lon: b.Lon}

	_, km := haversine.Distance(origin, dest)

	return km * 1000 // convert to meters
}
