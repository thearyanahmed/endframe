package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"
)

type RideServiceMock struct {
	mock.Mock
}

func (s *RideServiceMock) GetMinimumTripDistance() float64 {
	args := s.Called()
	return args.Get(0).(float64)
}

func (s *RideServiceMock) IsRideAvailable(ride locationEntity.Ride) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) GetRoute(origin, dest locationEntity.Coordinate) []locationEntity.Coordinate {
	args := s.Called()
	return args.Get(0).([]locationEntity.Coordinate)
}

func (s *RideServiceMock) RecordNewRideEvent(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error) {
	args := s.Called()
	return args.Get(0).(locationEntity.Event), args.Error(1)
}

func (s *RideServiceMock) RecordLocationUpdate(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error) {
	args := s.Called()
	return args.Get(0).(locationEntity.Event), args.Error(1)
}

func (s *RideServiceMock) DistanceIsGreaterThanMinimumDistance(origin, destination locationEntity.Coordinate) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) FindRideInLocation(ctx context.Context, rideUuid string, rideLocation locationEntity.Coordinate) (locationEntity.Ride, error) {
	args := s.Called()
	return args.Get(0).(locationEntity.Ride), args.Error(1)
}

func (s *RideServiceMock) UpdateRideLocation(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error) {
	args := s.Called()
	return args.Get(0).(locationEntity.Event), args.Error(1)
}

func (s *RideServiceMock) CanBeUpdatedViaRiderApp(ctx context.Context, rideUuid string, loc locationEntity.Coordinate) (bool, error) {
	args := s.Called()
	return args.Get(0).(bool), args.Error(1)
}

func (s *RideServiceMock) ResetMock() {
	s.Mock = mock.Mock{}
}
