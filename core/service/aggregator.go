package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/core/config"
	"github.com/thearyanahmed/nordsec/core/repository"
)

type ServiceAggregator struct {
	*RideService
}

func NewServiceAggregator(config *config.Specification, logger *log.Logger) (*ServiceAggregator, error) {
	redis, err := repository.NewRedisClient(config.GetRedisAddr(), config.GetRedisPassword())

	if err != nil {
		return &ServiceAggregator{}, err
	}

	rideRepo := repository.NewRideRepository(redis)
	rideSvc := NewRideService(rideRepo, logger)

	aggregator := &ServiceAggregator{
		RideService: rideSvc,
	}

	return aggregator, nil
}
