package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/core/config"
	"github.com/thearyanahmed/nordsec/core/repository"
	"github.com/thearyanahmed/nordsec/services/location"
)

type ServiceAggregator struct {
	*RideService
	LocationSvc *location.Service
}

func NewServiceAggregator(config *config.Specification, logger *log.Logger) (*ServiceAggregator, error) {
	// @todo extract to different layer
	redis, err := repository.NewRedisClient(config.GetRedisAddr(), config.GetRedisPassword())

	if err != nil {
		return &ServiceAggregator{}, err
	}

	rideRepo := repository.NewRideRepository(redis)
	rideSvc := NewRideService(rideRepo, logger)

	locSvc := location.NewLocationService(redis)

	aggregator := &ServiceAggregator{
		RideService: rideSvc,
		LocationSvc: locSvc,
	}

	return aggregator, nil
}
