package entity

import (
	"github.com/mmcloughlin/geohash"
)

// Area
// ┌──────────────────────────┐
// │(x4,y4)            (x3,y3)│
// │                          │
// │(x1,y1)            (x2,y2)│
// └──────────────────────────┘
// (X1Y1) : left bottom, (X2Y2) right bottom
// (X3Y3) : right top, (X4Y4) left top
type Area struct {
	X1Y1, X2Y2, X3Y3, X4Y4 Coordinate
}

// ToBoundingBox Create bounding box from an area
func (a *Area) ToBoundingBox() geohash.Box {
	return geohash.Box{
		MinLat: a.X1Y1.Lat,
		MaxLat: a.X3Y3.Lat,
		MinLng: a.X1Y1.Lon,
		MaxLng: a.X3Y3.Lon,
	}
}

// GetNeighbourGeohashFromCenter Get the neighbours based in the center of the bounding box's area
func (a *Area) GetNeighbourGeohashFromCenter(chars uint) []string {
	center := NewCoordinate(a.ToBoundingBox().Center())

	encodedCenter := geohash.EncodeWithPrecision(center.Lat, center.Lon, chars)

	return geohash.Neighbors(encodedCenter)
}
