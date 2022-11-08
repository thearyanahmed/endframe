package service

import (
	"github.com/thearyanahmed/nordsec/core/shared"
)

type rideRepository interface{}

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

func (s *RideService) UpdateRideStatus() {
	s.logger.Trace("update ride status invoked")
}
