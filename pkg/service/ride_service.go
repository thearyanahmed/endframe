package service

import (
	"context"
	"fmt"
	locationEntity "github.com/thearyanahmed/nordsec/services/location/entity"

	"github.com/thearyanahmed/nordsec/pkg/entity"
	"github.com/thearyanahmed/nordsec/pkg/repository"
	"github.com/thearyanahmed/nordsec/pkg/shared"
)

type rideRepository interface {
	UpdateLocation(ctx context.Context, uuid string, lat, long float64) (repository.RideLocationSchema, error)
	FindById(ctx context.Context, uuid string) (repository.RideLocationSchema, error)
}

type locationService interface {
	RecordRideEvent(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error)
}

type RideService struct {
	repository      rideRepository
	logger          shared.LoggerInterface
	locationService locationService
}

func NewRideService(repo rideRepository, locationSvc locationService, logger shared.LoggerInterface) *RideService {
	return &RideService{
		repository:      repo,
		logger:          logger,
		locationService: locationSvc,
	}
}

func (s *RideService) UpdateRideLocation(ctx context.Context, event locationEntity.Event) (locationEntity.Event, error) {
	rideEvent, err := s.locationService.RecordRideEvent(ctx, event)

	if err != nil {
		go func() {
			s.logger.Error(err)
		}()

		return locationEntity.Event{}, err
	}

	return rideEvent, nil
}

func (s *RideService) FindById(ctx context.Context, uuid string) (entity.RideLocationEntity, error) {
	panic("deprecated")
	loc, err := s.repository.FindById(ctx, uuid)

	if err != nil {
		return entity.RideLocationEntity{}, err
	}

	return loc.ToEntity(), nil
}

func (s *RideService) FindNearByRides(ctx context.Context, area locationEntity.Area) ([]entity.RideEntity, error) {
	panic("deprecated")
	fmt.Println("[ride service] GOT AREA", area)

	return []entity.RideEntity{}, nil
}
