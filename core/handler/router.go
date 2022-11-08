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

			// Update ride position
			r.Post("/update/ride-position", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("update ride position"))
			})

			

		})
	})

	return r
}
