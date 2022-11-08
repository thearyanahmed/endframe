package handler

import "net/http"

type activateRideHandler struct{}

func NewActiveRideHandler() *activateRideHandler {
	return &activateRideHandler{}
}

func (h *activateRideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}
