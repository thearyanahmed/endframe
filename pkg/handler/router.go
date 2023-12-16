package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/thearyanahmed/endframe/pkg/config"
	"github.com/thearyanahmed/endframe/pkg/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(conf *config.Specification, svcAggregator *service.ServiceAggregator, logger *log.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Return health check status of these services and pings redis, kafka and db(ping all related services)
			r.Get("/health-check", NewHealthCheckHandler(svcAggregator.GetKeyValueDataStore()).ServeHTTP)

			// This endpoint should return all available vehicles within an area also
			// should support filter by query params
			// r.With(apiMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
			// Get("/users/{id}", NewUserHandler(svcAggregator.UserSvc).ServeHTTP)
		})
	})

	return r
}
