package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/core/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// content type middleware?

			// @todo Add client authorization token middleware. Have different tokens for different routes.
			// use static tokens for simplicity

			// Return health check status of this services and pings redis, kafka and db? (ping all related services)
			// @q: ping user goroutines?
			r.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello world"))
			})

			// This endpoint should return all available vehichles within an area
			// also should support filter by query params
			// optional: perhaps we can also filter by the radius
			r.Get("/rides-near-by/{lat}/{long}/", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(fmt.Sprintf("lat: %s,long: %s", chi.URLParam(r, "lat"), chi.URLParam(r, "long"))))
			})

			// This endpoint spawns a ride from taking lat long and some other values from
			// the rider app.
			r.Post("/ride/spwan", func(w http.ResponseWriter, r *http.Request) {
				
			})

			// This endpoint takes input from the input, validates it.
			// Upon succesful validation, it creates creates 1 database entry, updates redis & go kafka().
			r.Post("/ride/start", func(w http.ResponseWriter, r *http.Request) {

			})

			// Update ride position takes a ride uuid, updates ride details on kafka.
			// But also need to to push to redis because users can query rides-near-by.
			// That query will be done from redis.
			r.Post("/notify/position", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("update ride position"))
			})

			// Validate the request, update to redis, kafka and database.
			r.Post("/ride/end", func(w http.ResponseWriter, r *http.Request) {

			})

			// Takes uuid, finds the details of the ride
			r.Get("/ride/{uuid}/view", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("get ride details"))
			})
		})
	})

	return r
}
