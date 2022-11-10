package main

import (
	"fmt"
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

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/mmcloughlin/geohash"
)

type Coord struct {
	Latitude  float64
	Longitude float64
}

func main() {
	twoPair := getTwoPair()
	fmt.Println(twoPair)

	box := makeBox(getCorners(twoPair))

	fmt.Println(box)

	centerX, centerY := box.Center()

	encodedCenter := geohash.EncodeWithPrecision(centerX, centerY, 6)
	fmt.Println(encodedCenter)

	neighbours := geohash.Neighbors(encodedCenter)

	fmt.Println("neighbours")
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

func latitudeRange() (float64, float64) {
	return 52.30000, 52.50000
}

func longitudeRange() (float64, float64) {
	return 13.46000, 13.54000
}

func spawnRiders() {
	if len(os.Args) < 4 {
		log.Fatal("missing required parameters. see readme.", os.Args)
	}

	endpoint := os.Args[3]

	requests, err := strconv.Atoi(os.Args[1])

	if err != nil {
		requests = 100
	}

	apiKey := os.Args[2]

	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	ch := make(chan string)

	var wg sync.WaitGroup

	wg.Add(requests)
	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := 0; i < requests; i++ {
		go makeRequest(i, ch, apiKey, endpoint, &wg)
	}

	defer func() {
		fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

// TODO organize
func makeRequest(i int, ch chan<- string, apiKey, endpoint string, wg *sync.WaitGroup) {
	defer wg.Done()

	form := url.Values{}

	form.Set("uuid", uuid.New().String())
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

	// todo take from arguments
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
