package middleware

import (
	"net/http"

	"github.com/thearyanahmed/nordsec/core/presenter"
)

func ValidateContentTypeMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			presenter.ErrorResponse(w, r, presenter.ErrInvalidContentType())
			return
		}
		next.ServeHTTP(w, r.WithContext(r.Context()))
	}
	return http.HandlerFunc(fn)
}
