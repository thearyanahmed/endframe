package location

import (
	"fmt"
	"github.com/mmcloughlin/geohash"
)

type Service struct {
	hashLength uint

	ds datastore
}

type datastore struct {
	ridesKey string
}

func NewGeoHash() Service {
	ds := datastore{
		ridesKey: "trips",
	}

	return Service{
		hashLength: uint(6),
		ds:         ds,
	}
}

func (s *Service) UpdateRideLocation(rideId string, location Coordinate) error {
	// get the location geohash
	ghash := location.Geohash(s.hashLength)

	err := geohash.Validate(ghash)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GetRidesIn(area Area) {

}
