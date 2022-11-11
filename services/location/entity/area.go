package entity

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
