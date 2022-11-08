package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	ch := make(chan string)

	var wg sync.WaitGroup

	requests := 100

	wg.Add(requests)
	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := 0; i < requests; i++ {
		go makeRequest(i, ch, &wg)
	}

	defer func() {
		fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

func makeRequest(i int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	ch <- fmt.Sprintf("[%d-start]:%s", i, "starting the trip")

	time.Sleep(1 * time.Second) // first start the event
	// the response should have some sort of [] of coordinates

	// imagine we have 10 points
	for x := 0; x < 10; x++ {
		ch <- fmt.Sprintf("[%d-x]:%s", i, "broadcasting event")
		time.Sleep(time.Second * 3)
	}

	ch <- fmt.Sprintf("[%d-result]:%s", i, "completing trip")
}
