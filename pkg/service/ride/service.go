package ride

import (
	"context"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"time"
)

type locationService interface {
	DistanceInMeters(origin, dest entity.Coordinate) float64
	GetRidesInArea(ctx context.Context, area entity.Area) ([]entity.Ride, error)
	RecordRideEvent(ctx context.Context, event entity.Event) (entity.Event, error)
	GetRoute(origin, destination entity.Coordinate, intervalPoints int) []entity.Coordinate
	FindRideInLocation(ctx context.Context, rideUuid string, origin entity.Coordinate) (entity.Ride, error)
	StartCooldownForRide(ctx context.Context, rideUuid string, timestamp int64, duration time.Duration) error

	UpdateRideEventCurrentStatus(ctx context.Context, event entity.Event) error
	GetRideEventByUuid(ctx context.Context, rideUuid string) (entity.Event, error)
}

type RideService struct {
	locationService       locationService
	minTripDistance       float64 // minimum distance between origin and destination, in meters
	inRouteIntervalPoints int
	cooldownMode          int64
}

func NewRideService(locationSvc locationService) *RideService {
	return &RideService{
		locationService:       locationSvc,
		minTripDistance:       500,
		inRouteIntervalPoints: 15, // how many points will be plotted between origin and destination (origin and destination are inclusive)
		cooldownMode:          3,  // in second
	}
}

func (s *RideService) getCooldownModeDurationInSeconds() time.Duration {
	return time.Duration(s.cooldownMode) * time.Second
}

func (s *RideService) GetMinimumTripDistance() float64 {
	return s.minTripDistance
}

func (s *RideService) RecordNewRideEvent(ctx context.Context, event entity.Event) (entity.Event, error) {
	event.SetStateAsInRoute().SetCurrentTimestamp().SetNewTripUuid()

	return s.UpdateRideLocation(ctx, event)
}

func (s *RideService) RecordLocationUpdate(ctx context.Context, event entity.Event) (entity.Event, error) {
	event.SetStateAsInRoute().SetCurrentTimestamp()

	return s.UpdateRideLocation(ctx, event)
}

func (s *RideService) RecordEndRideEvent(ctx context.Context, event entity.Event) (entity.Event, error) {
	event.SetStateAsRoaming().SetCurrentTimestamp()

	// @todo update ride:uuid current status

	return s.locationService.RecordRideEvent(ctx, event)
}

func (s *RideService) EnterCooldownMode(ctx context.Context, event entity.Event) error {
	return s.locationService.StartCooldownForRide(ctx, event.RideUuid, time.Now().Unix(), s.getCooldownModeDurationInSeconds())
}

func (s *RideService) UpdateRideLocation(ctx context.Context, event entity.Event) (entity.Event, error) {
	rideEvent, err := s.locationService.RecordRideEvent(ctx, event)

	if err != nil {
		return entity.Event{}, err
	}

	return rideEvent, nil
}

func (s *RideService) GetRoute(origin, dest entity.Coordinate) []entity.Coordinate {
	return s.locationService.GetRoute(origin, dest, s.inRouteIntervalPoints)
}

func (s *RideService) CanBeUpdatedViaRiderApp(ctx context.Context, rideUuid string, loc entity.Coordinate) (bool, error) {
	ride, err := s.locationService.FindRideInLocation(ctx, rideUuid, loc)

	if err != nil {
		return false, err
	}

	return s.IsRideAvailable(ride), nil
}

func (s *RideService) IsRideAvailable(ride entity.Ride) bool {
	return s.isAvailableByState(ride.State)
}

func (s *RideService) isAvailableByState(state string) bool {
	return state != entity.StateInCooldown && state != entity.StateInRoute
}

func (s *RideService) DistanceIsGreaterThanMinimumDistance(origin, dest entity.Coordinate) bool {
	return s.locationService.DistanceInMeters(origin, dest) > s.minTripDistance
}

func (s *RideService) FindRideInLocation(ctx context.Context, rideUuid string, rideLocation entity.Coordinate) (entity.Ride, error) {
	return s.locationService.FindRideInLocation(ctx, rideUuid, rideLocation)
}

func (s *RideService) FindNearByRides(ctx context.Context, area entity.Area, stateFilter string) ([]entity.Ride, error) {
	rides, err := s.locationService.GetRidesInArea(ctx, area)

	if err != nil {
		return []entity.Ride{}, err
	}

	if s.isValidFilter(stateFilter) {
		return s.filterByState(rides, stateFilter), nil
	}

	return rides, nil
}

func (s *RideService) isValidFilter(state string) bool {
	return state == entity.StateInCooldown || state == entity.StateRoaming || state == entity.StateInRoute
}

func (s *RideService) filterByState(collection []entity.Ride, state string) []entity.Ride {
	filtered := make([]entity.Ride, 0)

	for _, col := range collection {
		if col.State == state {
			filtered = append(filtered, col)
		}
	}

	return filtered
}

func (s *RideService) SetRideCurrentStatus(ctx context.Context, event entity.Event) error {
	return s.locationService.UpdateRideEventCurrentStatus(ctx, event)
}

func (s *RideService) GetRideEventByUuid(ctx context.Context, rideUuid string) (entity.Event, error) {
	return s.locationService.GetRideEventByUuid(ctx, rideUuid)
}

// TripHasEnded when a trip has ended, it should have any trip uuid and should be available
func (s *RideService) TripHasEnded(event entity.Event) bool {
	return event.TripUuid == "" && event.PassengerUuid == "" && (s.isRoaming(event.State) || s.isInCooldown(event.State))
}

func (s *RideService) isRoaming(state string) bool {
	return state == entity.StateRoaming
}

func (s *RideService) isInCooldown(state string) bool {
	return state == entity.StateInCooldown
}

func (s *RideService) IsInRoute(state string) bool {
	return state == entity.StateInRoute
}
