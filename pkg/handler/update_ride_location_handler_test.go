package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/thearyanahmed/nordsec/pkg/service"
	"github.com/thearyanahmed/nordsec/pkg/testutil"
)

type UpdateRideLocationHandlerTestSuite struct {
	suite.Suite

	usecase *service.RideServiceMock
}

func TestUpdateRideLocationHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateRideLocationHandlerTestSuite))
}

func (s *UpdateRideLocationHandlerTestSuite) SetupTest() {
	s.usecase = &service.RideServiceMock{}

	defer mock.AssertExpectationsForObjects(s.T(), s.usecase)
}

func (s *UpdateRideLocationHandlerTestSuite) TestRequestFailsWithoutValidFormData() {
	res := s.response(nil)

	assert.Equal(s.T(), http.StatusBadRequest, res.Code)
}

func (s *UpdateRideLocationHandlerTestSuite) response(data url.Values) *httptest.ResponseRecorder {
	// endpoint
	rr := testutil.NewAPIResponseBuilder().
		Method(http.MethodPost).
		URL("/api/v1/ride/activate").
		Body(data).
		Handler(NewUpdateRideLocationHandler(s.usecase).ServeHTTP).
		IsFormUrlEncoded().
		Build()

	return rr
}
