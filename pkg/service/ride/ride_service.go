package ride

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
)

type locationService interface {
	DistanceInMeters(origin, dest entity.Coordinate) float64
	GetRidesInArea(ctx context.Context, area entity.Area) ([]entity.Ride, error)
	RecordRideEvent(ctx context.Context, event entity.Event) (entity.Event, error)
	GetRoute(origin, destination entity.Coordinate, intervalPoints int) []entity.Coordinate
	FindRideInLocation(ctx context.Context, rideUuid string, origin entity.Coordinate) (entity.Ride, error)
}

type RideService struct {
	locationService       locationService
	minTripDistance       float64 // minimum distance between origin and destination, in meters
	inRouteIntervalPoints int
}

func NewRideService(locationSvc locationService) *RideService {
	return &RideService{
		locationService:       locationSvc,
		minTripDistance:       500,
		inRouteIntervalPoints: 15,
	}
}

func (s *RideService) GetMinimumTripDistance() float64 {
	return s.minTripDistance
}

func (s *RideService) RecordNewRideEvent(ctx context.Context, event entity.Event) (entity.Event, error) {
	event.SetStateAsInRoute().SetCurrentTimestamp().SetNewTripUuid()

	return s.UpdateRideLocation(ctx, event)
}

func (s *RideService) RecordLocationUpdate(ctx context.Context, event entity.Event) (entity.Event, error) {
	event.SetStateAsInRoute().SetCurrentTimestamp()

	return s.UpdateRideLocation(ctx, event)
}

func (s *RideService) UpdateRideLocation(ctx context.Context, event entity.Event) (entity.Event, error) {
	rideEvent, err := s.locationService.RecordRideEvent(ctx, event)

	if err != nil {
		return entity.Event{}, err
	}

	return rideEvent, nil
}

func (s *RideService) GetRoute(origin, dest entity.Coordinate) []entity.Coordinate {
	return s.locationService.GetRoute(origin, dest, s.inRouteIntervalPoints)
}

func (s *RideService) CanBeUpdatedViaRiderApp(ctx context.Context, rideUuid string, loc entity.Coordinate) (bool, error) {
	ride, err := s.locationService.FindRideInLocation(ctx, rideUuid, loc)

	if err != nil {
		return false, err
	}

	return s.IsRideAvailable(ride), nil
}

func (s *RideService) IsRideAvailable(ride entity.Ride) bool {
	return ride.State != entity.StateInCooldown && ride.State != entity.StateInRoute
}

func (s *RideService) DistanceIsGreaterThanMinimumDistance(origin, dest entity.Coordinate) bool {
	return s.locationService.DistanceInMeters(origin, dest) < s.minTripDistance
}

func (s *RideService) FindRideInLocation(ctx context.Context, rideUuid string, rideLocation entity.Coordinate) (entity.Ride, error) {
	return s.locationService.FindRideInLocation(ctx, rideUuid, rideLocation)
}

func (s *RideService) FindNearByRides(ctx context.Context, area entity.Area) ([]entity.Ride, error) {
	return s.locationService.GetRidesInArea(ctx, area)
}
