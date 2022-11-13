package main

type Coord struct {
	Latitude  float64
	Longitude float64
}

func main() {
	spawnRiders()

	// get near-by rides (concurrent)

	// start a random ride

	// ping locations

	// end ride

	// end of story
}

func latitudeRange() (float64, float64) {
	return 52.30000, 52.50000
}

func longitudeRange() (float64, float64) {
	return 13.46000, 13.54000
}
