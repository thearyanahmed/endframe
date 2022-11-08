package handler

import (
	"net/http"
	"time"

	"github.com/thearyanahmed/nordsec/core/presenter"
)

// @todo add dependend modules, eg: redis, kafka, db
type healthCheckHandler struct{}

func NewHealthCheckHandler() *healthCheckHandler {
	return &healthCheckHandler{}
}

// @todo update with pings
func (h *healthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := presenter.HealthCheck("ok", time.Now())

	presenter.RenderJsonResponse(w, r, http.StatusOK, res)
}
