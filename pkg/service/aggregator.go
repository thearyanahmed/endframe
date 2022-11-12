package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/pkg/config"
	"github.com/thearyanahmed/nordsec/pkg/repository"
	"github.com/thearyanahmed/nordsec/services/location"
)

type ServiceAggregator struct {
	*RideService
	LocationSvc *location.Service
}

func NewServiceAggregator(config *config.Specification, _ *log.Logger) (*ServiceAggregator, error) {
	// @todo extract to different layer
	redis, err := repository.NewRedisClient(config.GetRedisAddr(), config.GetRedisPassword())

	if err != nil {
		return &ServiceAggregator{}, err
	}

	locSvc := location.NewLocationService(redis, config.GetRedisLocationsKey())
	//rideRepo := repository.NewRideRepository(redis)

	rideSvc := NewRideService(locSvc)

	aggregator := &ServiceAggregator{
		RideService: rideSvc,
		LocationSvc: locSvc,
	}

	return aggregator, nil
}
