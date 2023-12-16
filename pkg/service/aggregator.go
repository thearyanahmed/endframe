package service

import (
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/endframe/pkg/config"
	"github.com/thearyanahmed/endframe/pkg/repository"
	"github.com/thearyanahmed/endframe/pkg/service/user"
)

type ServiceAggregator struct {
	UserSvc *user.UserService
	kvStore *redis.Client
}

func NewServiceAggregator(config *config.Specification, _ *log.Logger) (*ServiceAggregator, error) {
	kvStore, err := repository.NewRedisClient(config.GetRedisAddr(), config.GetRedisPassword())

	if err != nil {
		return &ServiceAggregator{}, err
	}

	userRepo := user.NewUserRepository()
	userSvc := user.NewUserService(userRepo)

	aggregator := &ServiceAggregator{
		UserSvc: userSvc,
		kvStore: kvStore,
	}

	return aggregator, nil
}

func (s *ServiceAggregator) GetKeyValueDataStore() *redis.Client {
	return s.kvStore
}
