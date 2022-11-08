package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type activateRideUsecase interface {
	UpdateRideStatus()
}

type activateRideHandler struct {
	usecase activateRideUsecase
	logger  *log.Logger
}

func NewActivateRideHandler(usecase activateRideUsecase, logger *log.Logger) *activateRideHandler {
	return &activateRideHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *activateRideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// validate the request
	// put to redis
	// put to kafka
	// return a response

	h.usecase.UpdateRideStatus()

	w.Write([]byte("done, check console log"))
}
