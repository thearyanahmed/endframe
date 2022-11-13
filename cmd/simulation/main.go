package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/mmcloughlin/geohash"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Coord struct {
	Latitude  float64
	Longitude float64
}

type Ride struct {
	RideUuid string  `json:"ride_uuid"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	State    string  `json:"state"`
}

type Route struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Trip struct {
	Message string  `json:"message"`
	Route   []Route `json:"route"`
	Event   struct {
		Uuid          string  `json:"uuid"`
		RideUuid      string  `json:"ride_uuid"`
		Lat           float64 `json:"lat"`
		Lon           float64 `json:"lon"`
		PassengerUuid string  `json:"passenger_uuid"`
		ClientUuid    string  `json:"client_uuid"`
		TripUuid      string  `json:"trip_uuid"`
		Timestamp     int     `json:"timestamp"`
		State         string  `json:"state"`
	} `json:"event"`
}

var rides []Ride

func main() {
	rand.Seed(time.Now().Unix())

	SpawnRiders()

	// get near-by rides (concurrent)
	GetNearByRides()

	// start a random ride
	SimulateTrip()
}

func SimulateTrip() {
	// select randomRide
	ch := make(chan string)
	numberOfTrips := 5
	var wg sync.WaitGroup

	wg.Add(numberOfTrips)
	go func() {
		wg.Wait()
		close(ch)
	}()

	startTripEndpoint := fmt.Sprintf("%s/trip/start", GetBaseUrl())

	for i := 0; i < numberOfTrips; i++ {
		go startTrip(ch, startTripEndpoint, GetClientApiKey(), randomTripData(), &wg)
	}

	for v := range ch {
		fmt.Println(v)
	}
}

func travelToDestination(trip Trip, wg *sync.WaitGroup) {
	defer wg.Done()

	endpoint := fmt.Sprintf("%s/trip/notify/location", GetBaseUrl())

	routeLen := len(trip.Route)

	i := 0

	for ; i < routeLen-1; i++ {
		r := trip.Route[i]

		fmt.Printf("Notifying server about trip: %s location: %v, %v\nWill need 3 seconds (time.Sleep) to travel to next destination.", trip.Event.TripUuid, r.Lat, r.Lon)
		pingLocation(trip, r, endpoint)

		time.Sleep(time.Second * 3)
	}

	reachDestination(trip, routeLen)
}

func reachDestination(trip Trip, routeLen int) {
	defer fmt.Printf("reached destination:%s\n", trip.Event.TripUuid)

	lastStop := trip.Route[routeLen-1]
	endTripEndpoint := fmt.Sprintf("%s/trip/end", GetBaseUrl())

	form := notifyTripLocationFormData(lastStop, trip.Event.TripUuid, trip.Event.RideUuid, trip.Event.PassengerUuid, trip.Event.ClientUuid)
	req, err := http.NewRequest("POST", endTripEndpoint, strings.NewReader(form.Encode()))

	req.Header.Add("Authorization", GetClientApiKey())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		fmt.Println("could not build request")
		return
	}

	fmt.Printf("reaching destination for tripid:%s\n", trip.Event.TripUuid)
	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
}

func pingLocation(trip Trip, r Route, endpoint string) bool {
	form := notifyTripLocationFormData(r, trip.Event.TripUuid, trip.Event.RideUuid, trip.Event.PassengerUuid, trip.Event.ClientUuid)
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))

	req.Header.Add("Authorization", GetClientApiKey())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return true
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	return false
}

func notifyTripLocationFormData(route Route, tripUuid, riderUuid, passengerUuid, clientUuid string) url.Values {
	r := url.Values{}

	r.Set("latitude", fmt.Sprintf("%.6f", route.Lat))
	r.Set("longitude", fmt.Sprintf("%.6f", route.Lon))
	r.Set("passenger_uuid", passengerUuid)
	r.Set("ride_uuid", riderUuid)
	r.Set("trip_uuid", tripUuid)
	r.Set("client_uuid", clientUuid)

	return r
}

func startTrip(ch chan<- string, endpoint, apiKey string, form url.Values, wg *sync.WaitGroup) {
	//defer wg.Done()
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))

	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		ch <- err.Error()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		ch <- err.Error()
		return
	}

	defer resp.Body.Close()

	var result Trip
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		ch <- err.Error()
		wg.Done()
		return
	}

	go travelToDestination(result, wg)
}

func randomTripData() url.Values {
	r := url.Values{}

	randomRide := rides[rand.Intn(len(rides))]

	r.Set("origin_latitude", fmt.Sprintf("%.6f", randomRide.Lat))
	r.Set("origin_longitude", fmt.Sprintf("%.6f", randomRide.Lon))
	r.Set("destination_latitude", fmt.Sprintf("%.6f", randomRide.Lon+0.3))
	r.Set("destination_longitude", fmt.Sprintf("%.6f", randomRide.Lon+0.3))
	r.Set("ride_uuid", randomRide.RideUuid)
	r.Set("client_uuid", uuid.New().String())

	return r
}

func latitudeRange() (float64, float64) {
	return 52.50000, 52.50000
}

func longitudeRange() (float64, float64) {
	return 13.46000, 13.54000
}

func GetBaseUrl() string {
	return os.Args[4]
}

func GetNumberOfRequests() int {
	if num, err := strconv.Atoi(os.Args[1]); err != nil {
		return 500
	} else {
		return num
	}
}

func GetRiderApiKey() string {
	return os.Args[2]
}

func GetClientApiKey() string {
	return os.Args[3]
}

func GetNearByRides() {
	endpoint := fmt.Sprintf("%s/rides/near-by", GetBaseUrl())
	apiKey := GetClientApiKey()
	requests := GetNumberOfRequests()

	var wg sync.WaitGroup

	wg.Add(requests)

	ch := make(chan string)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := 0; i < requests; i++ {
		go getRidesNearArea(ch, apiKey, endpoint, &wg)
	}

	for v := range ch {
		fmt.Println(v)
	}
}

func getRidesNearArea(ch chan<- string, apiKey, endpoint string, wg *sync.WaitGroup) {
	defer wg.Done()

	x1y1 := getRandomCoordinates()
	x3y3 := getRandomCoordinates()

	uri := fmt.Sprintf("%s?x1=%v&y1=%v&x3=%v&y3=%v", endpoint, x1y1.Latitude, x1y1.Longitude, x3y3.Latitude, x3y3.Longitude)

	req, err := http.NewRequest("GET", uri, strings.NewReader(""))

	req.Header.Add("Authorization", apiKey)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		ch <- err.Error()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		ch <- err.Error()
		return
	}

	defer resp.Body.Close()

	var result []Ride
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		ch <- err.Error()
		return
	}

	rides = append(rides, result...)

	ch <- fmt.Sprintf("in area bounded by x1y1:%v,%v and x3y3:%v,%v number of rides found : %d", x1y1.Latitude, x1y1.Longitude, x3y3.Latitude, x3y3.Longitude, len(result))
}

func SpawnRiders() {
	if len(os.Args) < 4 {
		log.Fatal("missing required parameters. see readme.", os.Args)
	}

	endpoint := fmt.Sprintf("%s/ride/activate", GetBaseUrl())
	requests := 500
	apiKey := GetRiderApiKey()

	ch := make(chan string)

	var wg sync.WaitGroup

	wg.Add(requests)
	go func() {
		wg.Wait()
		close(ch)
	}()

	fmt.Printf("Will attempt to spawn %d riders\n", requests)

	for i := 0; i < requests; i++ {
		go spawnNewRider(ch, apiKey, endpoint, &wg)
	}

	for v := range ch {
		fmt.Println(v)
	}
}

func spawnNewRider(ch chan<- string, apiKey, endpoint string, wg *sync.WaitGroup) {
	defer wg.Done()

	form := url.Values{}

	form.Set("ride_uuid", uuid.New().String())
	lat, err := gofakeit.LatitudeInRange(latitudeRange())
	if err != nil {
		ch <- err.Error()
		return
	}

	lon, err := gofakeit.LatitudeInRange(longitudeRange())
	if err != nil {
		ch <- err.Error()
		return
	}

	form.Set("latitude", fmt.Sprintf("%f", lat))
	form.Set("longitude", fmt.Sprintf("%f", lon))

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))

	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		ch <- err.Error()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		ch <- err.Error()
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	ch <- string(body)
}

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
