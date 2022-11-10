package location

import (
	"github.com/mmcloughlin/geohash"
)

// Area
// X1Y1 is left bottom
// X2Y2 is right bottom
// X3Y3 is right top
// X4Y4 is left top
type Area struct {
	X1Y1, X2Y2, X3Y3, X4Y4 Coordinate
}

type Coordinate struct {
	Lat, Lon float64
}

func NewCoordinate(lat, lon float64) *Coordinate {
	return &Coordinate{
		Lat: lat,
		Lon: lon,
	}
}

func (a *Area) ToBoundingBox() geohash.Box {
	return geohash.Box{
		MinLat: a.X1Y1.Lat,
		MaxLat: a.X3Y3.Lat,
		MinLng: a.X1Y1.Lon,
		MaxLng: a.X3Y3.Lon,
	}
}

func (a *Area) GetNeighbourGeohashFromCenter(chars uint) []string {
	center := NewCoordinate(a.ToBoundingBox().Center())

	encodedCenter := geohash.EncodeWithPrecision(center.Lat, center.Lon, chars)

	return geohash.Neighbors(encodedCenter)
}

type RideEvent struct {
	Uuid          string  `json:"uuid"`
	RideUuid      string  `json:"ride_uuid"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	PassengerUuid string  `json:"passenger_uuid"`
	Timestamp     int64   `json:"timestamp"`
	State         string  `json:"state"` // in route, roaming
}
