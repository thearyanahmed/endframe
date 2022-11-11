package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/thearyanahmed/nordsec/pkg/config"
	apiMiddleware "github.com/thearyanahmed/nordsec/pkg/middleware"
	"github.com/thearyanahmed/nordsec/pkg/service"

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
			r.Get("/health-check", NewHealthCheckHandler().ServeHTTP)

			// This endpoint should return all available vehicles within an area also
			// should support filter by query params
			r.With(apiMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Get("/rides/near-by", NewNearByRidesHandler(svcAggregator.RideService, svcAggregator.LocationSvc).ServeHTTP)

			// This simulates the idea of a rider, who just came online
			r.With(apiMiddleware.ValidateContentTypeMiddleware).
				With(apiMiddleware.NewAuthorizeRiderMiddleware(conf.RiderApiKey, logger).Handle).
				Post("/ride/activate", NewUpdateRideLocationHandler(svcAggregator.RideService).ServeHTTP)

			// Start a trip.
			r.With(apiMiddleware.ValidateContentTypeMiddleware).
				With(apiMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Post("/trip/start", NewStartTripHandler(svcAggregator.LocationSvc).ServeHttp)

			// Update location while on trip
			r.With(apiMiddleware.ValidateContentTypeMiddleware).
				With(apiMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Post("/trip/notify/location", NewNotifyPositionHandler(svcAggregator.LocationSvc).ServeHTTP)

			// Notify when trip has ended
			r.With(apiMiddleware.ValidateContentTypeMiddleware).
				With(apiMiddleware.NewAuthorizeClientMiddleware(conf.ClientApiKey, logger).Handle).
				Post("/trip/end", NewEndTripHandler(svcAggregator.LocationSvc).ServeHTTP)
		})
	})

	return r
}
