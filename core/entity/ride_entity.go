package entity

type RideLocationEntity struct {
	UUID      string
	Latitude  float64
	Longitude float64
}

type RideEntity struct {
	UUID      string
	Latitude  float64
	Longitude float64
	Status    string // in route, available, cooldown
}
