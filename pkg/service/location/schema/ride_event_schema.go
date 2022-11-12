package schema

import (
	"github.com/google/uuid"
	entity2 "github.com/thearyanahmed/nordsec/pkg/service/location/entity"
)

// @todo BIG REFACTOR

type RideEventSchema struct {
	Uuid          uuid.UUID `json:"uuid"`
	RideUuid      string    `json:"ride_uuid"`
	Lat           float64   `json:"lat"`
	Lon           float64   `json:"lon"`
	PassengerUuid string    `json:"passenger_uuid"`
	TripUuid      string    `json:"trip_uuid"`
	Timestamp     int64     `json:"timestamp"`
	State         string    `json:"state"` // in route, roaming
}

func (s *RideEventSchema) ToEntity() entity2.Event {
	return entity2.Event{
		Uuid:          s.Uuid.String(),
		RideUuid:      s.RideUuid,
		Lat:           s.Lat,
		Lon:           s.Lon,
		TripUuid:      s.TripUuid,
		PassengerUuid: s.PassengerUuid,
		Timestamp:     s.Timestamp,
		State:         s.State,
	}
}

func (s *RideEventSchema) WithNewUuid() *RideEventSchema {
	s.Uuid = uuid.New()

	return s
}

func (s *RideEventSchema) ToRideEntity() entity2.Ride {
	return entity2.Ride{
		RideUuid: s.RideUuid,
		Lat:      s.Lat,
		Lon:      s.Lon,
		State:    s.State,
	}
}
func FromRideEventEntity(e entity2.Event) *RideEventSchema {
	s := &RideEventSchema{
		RideUuid:      e.RideUuid,
		Lat:           e.Lat,
		Lon:           e.Lon,
		PassengerUuid: e.PassengerUuid,
		TripUuid:      e.TripUuid,
		Timestamp:     e.Timestamp,
		State:         e.State,
	}

	if e.Uuid != "" {
		s.Uuid = uuid.Must(uuid.FromBytes([]byte(e.Uuid)))
	}

	return s
}

func FromRideEventCollectionToEntity(list map[string]RideEventSchema) []entity2.Ride {
	var rides []entity2.Ride

	for _, e := range list {
		rides = append(rides, e.ToRideEntity())
	}

	return rides
}
