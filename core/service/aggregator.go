package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/core/config"
	"github.com/thearyanahmed/nordsec/core/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceAggregator struct {
	*RideService
}

func NewServiceAggregator(config *config.Specification, db *mongo.Database, logger *log.Logger) (*ServiceAggregator, error) {
	// redis, err := repository.NewRedisClient(config.GetRedisAddr(), config.GetRedisPassword())

	// if err != nil {
	// 	return &ServiceAggregator{}, err
	// }

	rideRepo := repository.NewRideRepository(db.Collection("rides_db_test"))
	rideSvc := NewRideService(rideRepo, logger)

	aggregator := &ServiceAggregator{
		RideService: rideSvc,
	}

	return aggregator, nil
}
