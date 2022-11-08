package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/core/presenter"
)

type authorizeRiderMiddleware struct {
	authToken string
	logger    *log.Logger
}

func NewAuthorizeRiderMiddleware(token string, logger *log.Logger) *authorizeRiderMiddleware {
	return &authorizeRiderMiddleware{
		authToken: token,
		logger:    logger,
	}
}

func (m *authorizeRiderMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			go func() {
				m.logger.
					WithField("req-id", middleware.GetReqID(r.Context())).
					Error("attempt to authorize without auth token")
			}()
			presenter.ErrorResponse(w, r, presenter.ErrUnauthorized())
			return
		}

		if token != m.authToken {
			go func() {
				m.logger.
					WithField("req-id", middleware.GetReqID(r.Context())).
					WithField("used-token", token).
					Error("attempt to authorize without auth token")

			}()
			presenter.ErrorResponse(w, r, presenter.ErrUnauthorized())
			return
		}
		go func() {
			m.logger.
				WithField("req-id", middleware.GetReqID(r.Context())).
				Trace("authorized to rider app")
		}()
		next.ServeHTTP(w, r)
	})
}
