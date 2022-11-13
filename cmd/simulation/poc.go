package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mmcloughlin/geohash"
)

func poc() {
	twoPair := getTwoPair()

	box := makeBox(getCorners(twoPair))

	centerX, centerY := box.Center()

	encodedCenter := geohash.EncodeWithPrecision(centerX, centerY, 6)

	neighbours := geohash.Neighbors(encodedCenter)

	for _, v := range neighbours {
		lat, lon := geohash.Decode(v)
		fmt.Printf("lat: %.5f lon: %.5f geo: %v\n", lat, lon, v)
	}
}

func getCorners(dataset []Coord) (float64, float64, float64, float64) {
	minX, minY := dataset[0].Latitude, dataset[0].Longitude

	maxX, maxY := dataset[2].Latitude, dataset[3].Longitude

	return minX, minY, maxX, maxY
}

func makeBox(minX, minY, maxX, maxY float64) geohash.Box {
	return geohash.Box{
		MinLat: minX,
		MaxLat: maxX,
		MinLng: minY,
		MaxLng: maxY,
	}
}

func getTwoPair() []Coord {
	var cord []Coord

	// (0,0) bottom left
	cord = append(cord, Coord{
		Latitude:  latInRange(52.10000, 52.30000),
		Longitude: lonInRange(13.10000, 13.30000),
	})

	// (6,0)
	cord = append(cord, Coord{
		Latitude:  latInRange(52.10000, 52.30000),
		Longitude: lonInRange(13.10000, 13.30000),
	})

	// (6,6) // top right
	cord = append(cord, Coord{
		Latitude:  latInRange(52.31000, 52.50000),
		Longitude: lonInRange(13.40000, 13.60000),
	})

	cord = append(cord, Coord{
		Latitude:  latInRange(52.10000, 52.30000),
		Longitude: lonInRange(13.60000, 13.90000),
	})

	return cord
}

func latInRange(a, b float64) float64 {
	c, _ := gofakeit.LatitudeInRange(a, b)

	return c
}

func lonInRange(a, b float64) float64 {
	c, _ := gofakeit.LongitudeInRange(a, b)

	return c
}

func getRandomCoordinates() Coord {
	lat, _ := gofakeit.LatitudeInRange(latitudeRange())
	lon, _ := gofakeit.LongitudeInRange(longitudeRange())

	return Coord{lat, lon}
}
