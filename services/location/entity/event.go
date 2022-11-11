package entity

import (
	"time"
)

const StateRoaming = "roaming"
const StateInRoute = "in_route"
const StateInCooldown = "cooldown"

type Event struct {
	Uuid          string  `json:"uuid"`
	RideUuid      string  `json:"ride_uuid"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	PassengerUuid string  `json:"passenger_uuid"`
	TripUuid      string  `json:"trip_uuid"`
	Timestamp     int64   `json:"timestamp"`
	State         string  `json:"state"` // in route, roaming
}

func (r *Event) SetCurrentTimestamp() *Event {
	r.Timestamp = time.Now().Unix()

	return r
}

func (r *Event) SetStateAsRoaming() *Event {
	r.State = StateRoaming

	return r
}

func (r *Event) SetStateAsInRoute() *Event {
	r.State = StateInRoute

	return r
}

func (r *Event) SetStateAsInCooldown() *Event {
	r.State = StateInCooldown

	return r
}
