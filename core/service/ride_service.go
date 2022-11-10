package service

import (
	"context"
	"fmt"
	"github.com/thearyanahmed/nordsec/services/location"

	"github.com/thearyanahmed/nordsec/core/entity"
	"github.com/thearyanahmed/nordsec/core/repository"
	"github.com/thearyanahmed/nordsec/core/shared"
)

type rideRepository interface {
	UpdateLocation(ctx context.Context, uuid string, lat, long float64) (repository.RideLocationSchema, error)
	FindById(ctx context.Context, uuid string) (repository.RideLocationSchema, error)
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

func (s *RideService) UpdateRideLocation(ctx context.Context, uuid string, lat, long float64) (entity.RideLocationEntity, error) {
	panic("deprecated")
	schema, err := s.repository.UpdateLocation(ctx, uuid, lat, long)

	if err != nil {
		return entity.RideLocationEntity{}, err
	}

	// @todo add logging
	return schema.ToEntity(), nil
}

func (s *RideService) FindById(ctx context.Context, uuid string) (entity.RideLocationEntity, error) {
	loc, err := s.repository.FindById(ctx, uuid)

	if err != nil {
		return entity.RideLocationEntity{}, err
	}

	return loc.ToEntity(), nil
}

func (s *RideService) FindNearByRides(ctx context.Context, area location.Area) ([]entity.RideEntity, error) {
	fmt.Println("[ride service] GOT AREA", area)

	return []entity.RideEntity{}, nil
}
