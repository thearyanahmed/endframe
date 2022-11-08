package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	for i := 0; i < 100; i++ {
		go makeRequest(ch)
	}
	for i := 0; i < 100; i++ {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func makeRequest(ch chan<- string) {
	resp, _ := http.Get("https://jsonplaceholder.typicode.com/posts")
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- fmt.Sprintf(string(body))
}
