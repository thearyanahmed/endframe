package service

import (
	"context"
	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"
)

type locationService interface {
	DistanceInMeters(origin, dest locationEntity.Coordinate) float64
	GetRidesInArea(ctx context.Context, area locationEntity.Area) ([]locationEntity.Ride, error)
	RecordRideEvent(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error)
	GetRoute(origin, destination locationEntity.Coordinate, intervalPoints int) []locationEntity.Coordinate
	FindRideInLocation(ctx context.Context, rideUuid string, origin locationEntity.Coordinate) (locationEntity.Ride, error)
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

func (s *RideService) RecordRideEvent(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error) {
	event.SetStateAsInRoute().SetCurrentTimestamp().SetNewTripUuid()

	return s.UpdateRideLocation(ctx, event)
}

func (s *RideService) GetMinimumTripDistance() float64 {
	return s.minTripDistance
}

func (s *RideService) UpdateRideLocation(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error) {
	rideEvent, err := s.locationService.RecordRideEvent(ctx, event)

	if err != nil {
		return locationEntity.Event{}, err
	}

	return rideEvent, nil
}

func (s *RideService) GetRoute(origin, dest locationEntity.Coordinate) []locationEntity.Coordinate {
	return s.locationService.GetRoute(origin, dest, s.inRouteIntervalPoints)
}

func (s *RideService) CanBeUpdatedViaRiderApp(ctx context.Context, rideUuid string, loc locationEntity.Coordinate) (bool, error) {
	ride, err := s.locationService.FindRideInLocation(ctx, rideUuid, loc)

	if err != nil {
		return false, err
	}

	return s.IsRideAvailable(ride), nil
}

func (s *RideService) IsRideAvailable(ride locationEntity.Ride) bool {
	return ride.State != locationEntity.StateInCooldown && ride.State != locationEntity.StateInRoute
}

func (s *RideService) DistanceIsGreaterThanMinimumDistance(origin, dest locationEntity.Coordinate) bool {
	return s.locationService.DistanceInMeters(origin, dest) < s.minTripDistance
}

func (s *RideService) FindRideInLocation(ctx context.Context, rideUuid string, rideLocation locationEntity.Coordinate) (locationEntity.Ride, error) {
	return s.locationService.FindRideInLocation(ctx, rideUuid, rideLocation)
}

func (s *RideService) FindNearByRides(ctx context.Context, area locationEntity.Area) ([]locationEntity.Ride, error) {
	return s.locationService.GetRidesInArea(ctx, area)
}
