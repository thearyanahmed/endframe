package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/thearyanahmed/nordsec/core/config"
	coreMiddleware "github.com/thearyanahmed/nordsec/core/middleware"
	"github.com/thearyanahmed/nordsec/core/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(conf *config.Specification, svcAggregator *service.ServiceAggregator, logger *log.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/core/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Return health check status of these services and pings redis, kafka and db(ping all related services)
			r.Get("/health-check", NewHealthCheckHandler().ServeHTTP)

			// @todo Add client authorization token middleware. Have different tokens for different routes.
			// use static tokens for simplicity
			// This endpoint should return all available vehicles within an area
			// also should support filter by query params
			// optional: perhaps we can also filter by the radius
			r.With(coreMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Get("/rides/near-by", NewNearByRidesHandler(svcAggregator.RideService, svcAggregator.LocationSvc).ServeHTTP)

			// This endpoint spawns a ride from taking lat long and some other values from
			// the rider app.
			// This simulates the idea of a rider, who just came online
			r.With(coreMiddleware.ValidateContentTypeMiddleware).
				With(coreMiddleware.NewAuthorizeRiderMiddleware(conf.RiderApiKey, logger).Handle).
				Post("/ride/activate", NewUpdateRideLocationHandler(svcAggregator.RideService, svcAggregator.LocationSvc).ServeHTTP)

			// This endpoint takes input from the input, validates it.
			r.With(coreMiddleware.ValidateContentTypeMiddleware).
				With(coreMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Post("/trip/start", NewStartTripHandler(svcAggregator.LocationSvc).ServeHttp)

			// Update ride position takes a ride uuid, updates ride details on kafka.
			// But also need to push to redis because users can query rides-near-by.
			// That query will be done from redis.
			r.With(coreMiddleware.ValidateContentTypeMiddleware).
				With(coreMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Post("/trip/notify/location", NewNotifyPositionHandler(svcAggregator.LocationSvc).ServeHTTP)

			// Validate the request, update to redis, kafka and database.
			r.With(coreMiddleware.ValidateContentTypeMiddleware).
				With(coreMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Post("/trip/end", NewEndTripHandler(svcAggregator.LocationSvc).ServeHTTP)
		})
	})

	return r
}
