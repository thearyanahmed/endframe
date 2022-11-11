package presenter

import (
	"fmt"
	"net/http"
)

func ErrDistanceTooLowResponse(minDistance float64) *Response {
	return &Response{
		HttpStatusCode: http.StatusBadRequest,
		Message:        fmt.Sprintf("distance too low. minimum distance required %.2f meters or greater", minDistance),
	}
}

func ErrRideUnavailableResponse() *Response {
	return &Response{
		HttpStatusCode: http.StatusUnprocessableEntity,
		Message:        "ride is unavailable",
	}
}
