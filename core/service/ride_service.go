package service

import (
	"context"

	"github.com/thearyanahmed/nordsec/core/shared"
)

type rideRepository interface {
	UpdateLocation(ctx context.Context, uuid string, lat, long float64) error
}

type RideService struct {
	repository rideRepository
	logger     shared.LoggerInterface
}

func NewRideService(repo rideRepository, logger shared.LoggerInterface) *RideService {
	return &RideService{
		repository: repo,
		logger:     logger,
	}
}

func (s *RideService) UpdateRideLocation(ctx context.Context, uuid string, lat, long float64) error {
	err := s.repository.UpdateLocation(ctx, uuid, lat, long)

	if err != nil {
		return err
	}

	// @todo add logging
	return nil
}
