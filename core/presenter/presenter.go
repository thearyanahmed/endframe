package presenter

import (
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

type Response struct {
	HttpStatusCode int        `json:"-"`
	Message        string     `json:"message"`
	Details        url.Values `json:"details,omitempty"`
}

func (_ *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ErrUnauthorized() *Response {
	return &Response{
		HttpStatusCode: http.StatusUnauthorized,
		Message:        "unauthorized.",
	}
}

func FromErr(err error) *Response {
	return &Response{
		HttpStatusCode: http.StatusUnprocessableEntity, // @todo update here
		Message:        err.Error(),
	}
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, er *Response) {
	RenderJsonResponse(w, r, er.HttpStatusCode, er)
}

func RenderJsonResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	render.Status(r, statusCode)
	render.JSON(w, r, data)
}

func ErrorValidationFailed(validationErrors url.Values) *Response {
	return &Response{
		HttpStatusCode: http.StatusBadRequest,
		Message:        "validation failed",
		Details:        validationErrors,
	}
}

func ErrInvalidContentType() *Response {
	return &Response{
		HttpStatusCode: http.StatusNotAcceptable,
		Message:        "content-type must be application/x-www-form-urlencoded",
	}
}
