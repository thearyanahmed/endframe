package ride

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
)

type RideServiceMock struct {
	mock.Mock
}

func (s *RideServiceMock) RecordEndRideEvent(ctx context.Context, event entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) EnterCooldownMode(ctx context.Context, event entity.Event) error {
	args := s.Called()
	return args.Get(0).(error)
}

func (s *RideServiceMock) FindNearByRides(_ context.Context, _ entity.Area, _ string) ([]entity.Ride, error) {
	args := s.Called()
	return args.Get(0).([]entity.Ride), args.Error(1)
}

func (s *RideServiceMock) GetMinimumTripDistance() float64 {
	args := s.Called()
	return args.Get(0).(float64)
}

func (s *RideServiceMock) IsRideAvailable(ride entity.Ride) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) GetRoute(origin, dest entity.Coordinate) []entity.Coordinate {
	args := s.Called()
	return args.Get(0).([]entity.Coordinate)
}

func (s *RideServiceMock) RecordNewRideEvent(ctx context.Context, event entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) RecordLocationUpdate(ctx context.Context, event entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) DistanceIsGreaterThanMinimumDistance(origin, destination entity.Coordinate) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) FindRideInLocation(ctx context.Context, rideUuid string, rideLocation entity.Coordinate) (entity.Ride, error) {
	args := s.Called()
	return args.Get(0).(entity.Ride), args.Error(1)
}

func (s *RideServiceMock) UpdateRideLocation(ctx context.Context, event entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) CanBeUpdatedViaRiderApp(ctx context.Context, rideUuid string, loc entity.Coordinate) (bool, error) {
	args := s.Called()
	return args.Get(0).(bool), args.Error(1)
}

func (s *RideServiceMock) ResetMock() {
	s.Mock = mock.Mock{}
}
