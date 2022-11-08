package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() {
	r := chi.NewRouter()

	r.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	fmt.Println("serving")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}
