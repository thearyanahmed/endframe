package entity

type RideLocationEntity struct {
	RideUuid  string
	Latitude  float64
	Longitude float64
}

type RideEntity struct {
	RideUuid  string
	Latitude  float64
	Longitude float64
	State     string // in route, available, cooldown
}
