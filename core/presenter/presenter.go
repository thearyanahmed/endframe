package presenter

import (
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

type Response struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Details url.Values `json:"details,omitempty"`
}

func (_ *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, er *Response) {
	render.Status(r, er.Code)
	RenderJsonResponse(w, r, er.Code, er)
}

func RenderJsonResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	render.Status(r, statusCode)
	render.JSON(w, r, data)
}

func ErrInvalidContentType() *Response {
	return &Response{
		Code:    http.StatusNotAcceptable,
		Message: "content-type must be application/x-www-form-urlencoded",
	}
}
