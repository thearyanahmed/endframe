package handler

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/hellofresh/health-go/v4"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
	"net/http"
	"time"
)

type (
	healthChecker interface {
		Ping(ctx context.Context) *redis.StatusCmd
	}

	healthCheckHandler struct {
		health *health.Health
	}
)

func NewHealthCheckHandler(datastore healthChecker) *healthCheckHandler {
	h, _ := health.New(
		health.WithChecks(
			health.Config{
				Name:    "redis", // @improvement should be configurable
				Timeout: 2 * time.Second,
				Check: func(context.Context) error {
					_, err := datastore.Ping(context.TODO()).Result()

					return err
				},
			},
		),
		health.WithComponent(health.Component{
			Name:    "api",
			Version: "v1",
		}))
	return &healthCheckHandler{h}
}

func (h *healthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.health == nil {
		presenter.RenderJsonResponse(w, r, http.StatusOK, map[string]interface{}{
			"status":    "unknown",
			"timestamp": time.Now(),
			"failures": map[string]interface{}{
				"health_check": "failed to initiate health checker",
			},
		})
		return
	}

	res := presenter.HealthCheck("ok", time.Now())

	presenter.RenderJsonResponse(w, r, http.StatusOK, res)
}
