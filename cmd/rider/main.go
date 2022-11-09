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
)

func main() {
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

func latitudeRange() (float64, float64) {
	return 52.30000, 52.50000
}

func longitudeRange() (float64, float64) {
	return 13.46000, 13.54000
}
