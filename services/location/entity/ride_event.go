package entity

type RideEvent struct {
	Uuid          string  `json:"uuid"`
	RideUuid      string  `json:"ride_uuid"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	PassengerUuid string  `json:"passenger_uuid"`
	TripUuid      string  `json:"trip_uuid"`
	Timestamp     int64   `json:"timestamp"`
	State         string  `json:"state"` // in route, roaming
}
