package service

import (
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/pkg/config"
	"github.com/thearyanahmed/nordsec/pkg/repository"
	"github.com/thearyanahmed/nordsec/pkg/service/location"
	"github.com/thearyanahmed/nordsec/pkg/service/ride"
)

type ServiceAggregator struct {
	*ride.RideService
	keyValueDatastore *redis.Client
	LocationSvc       *location.Service
}

func NewServiceAggregator(config *config.Specification, _ *log.Logger) (*ServiceAggregator, error) {
	datastore, err := repository.NewRedisClient(config.GetRedisAddr(), config.GetRedisPassword())

	if err != nil {
		return &ServiceAggregator{}, err
	}

	locationService := location.NewLocationService(datastore, config.GetRedisLocationsKey())

	rideSvc := ride.NewRideService(locationService, config.GetMinimumTripDistance(), config.GetCooldownDuration())

	aggregator := &ServiceAggregator{
		RideService:       rideSvc,
		LocationSvc:       locationService,
		keyValueDatastore: datastore,
	}

	return aggregator, nil
}

func (s *ServiceAggregator) GetKeyValueDataStore() *redis.Client {
	return s.keyValueDatastore
}
