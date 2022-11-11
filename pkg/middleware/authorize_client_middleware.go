package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
)

type authorizeClientMiddleware struct {
	authToken string
	logger    *log.Logger
}

func NewAuthorizeClientMiddleware(token string, logger *log.Logger) *authorizeClientMiddleware {
	return &authorizeClientMiddleware{
		authToken: token,
		logger:    logger,
	}
}

func (m *authorizeClientMiddleware) Handle(next http.Handler) http.Handler {
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
				Trace("authorized to client app")
		}()
		next.ServeHTTP(w, r)
	})
}
