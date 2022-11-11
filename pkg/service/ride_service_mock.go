package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"
)

type RideServiceMock struct {
	mock.Mock
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
