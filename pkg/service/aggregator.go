package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/pkg/config"
	"github.com/thearyanahmed/nordsec/pkg/repository"
	"github.com/thearyanahmed/nordsec/pkg/service/location"
	"github.com/thearyanahmed/nordsec/pkg/service/ride"
)

type ServiceAggregator struct {
	*ride.RideService
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

	rideSvc := ride.NewRideService(locSvc)

	aggregator := &ServiceAggregator{
		RideService: rideSvc,
		LocationSvc: locSvc,
	}

	return aggregator, nil
}
