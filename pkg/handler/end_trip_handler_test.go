package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/thearyanahmed/nordsec/pkg/service/ride"
	"github.com/thearyanahmed/nordsec/pkg/testutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type endTripHandlerTestSuite struct {
	suite.Suite
	rideService *ride.RideServiceMock
}

func (h *endTripHandlerTestSuite) SetupTest() {
	h.rideService = &ride.RideServiceMock{}

	defer mock.AssertExpectationsForObjects(h.T(), h.rideService)
}

func TestEndTripHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(endTripHandlerTestSuite))
}

func (h *endTripHandlerTestSuite) TestSanityCheck() {
	assert.Equal(h.T(), 1, 1)
}

func (h *endTripHandlerTestSuite) response(data url.Values) *httptest.ResponseRecorder {
	rr := testutil.NewAPIResponseBuilder().
		Method(http.MethodPost).
		URL("/api/v1/trip/end").
		Body(data).
		Handler(NewEndTripHandler(h.rideService).ServeHTTP).
		IsFormUrlEncoded().
		Build()

	return rr
}
