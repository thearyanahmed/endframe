package ride

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thearyanahmed/endframe/pkg/service/location/entity"
)

type RideServiceMock struct {
	mock.Mock
}

func (s *RideServiceMock) IsInRoute(state string) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) SetRideCurrentStatus(ctx context.Context, event entity.Event) error {
	args := s.Called()
	return args.Error(0)
}

func (s *RideServiceMock) GetRideEventByUuid(_ context.Context, _ string) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) TripHasEnded(_ entity.Event) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) RecordEndRideEvent(_ context.Context, _ entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) EnterCooldownMode(_ context.Context, _ entity.Event) error {
	args := s.Called()
	return args.Error(0)
}

func (s *RideServiceMock) FindNearByRides(_ context.Context, _ entity.Area, _ string) ([]entity.Ride, error) {
	args := s.Called()
	return args.Get(0).([]entity.Ride), args.Error(1)
}

func (s *RideServiceMock) GetMinimumTripDistance() float64 {
	args := s.Called()
	return args.Get(0).(float64)
}

func (s *RideServiceMock) IsRideAvailable(_ entity.Ride) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) GetRoute(_, _ entity.Coordinate) []entity.Coordinate {
	args := s.Called()
	return args.Get(0).([]entity.Coordinate)
}

func (s *RideServiceMock) RecordNewRideEvent(_ context.Context, _ entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) RecordLocationUpdate(_ context.Context, _ entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) DistanceIsGreaterThanMinimumDistance(_, _ entity.Coordinate) bool {
	args := s.Called()
	return args.Get(0).(bool)
}

func (s *RideServiceMock) FindRideInLocation(_ context.Context, _ string, _ entity.Coordinate) (entity.Ride, error) {
	args := s.Called()
	return args.Get(0).(entity.Ride), args.Error(1)
}

func (s *RideServiceMock) UpdateRideLocation(_ context.Context, _ entity.Event) (entity.Event, error) {
	args := s.Called()
	return args.Get(0).(entity.Event), args.Error(1)
}

func (s *RideServiceMock) CanBeUpdatedViaRiderApp(_ context.Context, _ string, _ entity.Coordinate) (bool, error) {
	args := s.Called()
	return args.Get(0).(bool), args.Error(1)
}

func (s *RideServiceMock) ResetMock() {
	s.Mock = mock.Mock{}
}
