package entity

type Ride struct {
	RideUuid string  `json:"ride_uuid"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	State    string  `json:"state"` // in route, roaming
}
