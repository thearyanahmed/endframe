package location

import (
	"github.com/google/uuid"
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

func (s *RideEventSchema) ToEntity() RideEvent {
	return RideEvent{
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

func fromRideEventEntity(e RideEvent) *RideEventSchema {
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

func fromRideEventCollection(collection []RideEventSchema) []RideEvent {
	var eventSchema []RideEvent

	for _, e := range collection {
		eventSchema = append(eventSchema, e.ToEntity())
	}

	return eventSchema
}

func (s *RideEventSchema) WithNewUuid() *RideEventSchema {
	s.Uuid = uuid.New()

	return s
}

func (s *RideEventSchema) ToRideEntity() Ride {
	return Ride{
		RideUuid: s.RideUuid,
		Lat:      s.Lat,
		Lon:      s.Lon,
		State:    s.State,
	}
}
