package handler

import (
	"encoding/json"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/thearyanahmed/nordsec/pkg/service"
	"github.com/thearyanahmed/nordsec/pkg/testutil"
)

type startTripHandlerTestSuite struct {
	suite.Suite
	minTripDistance float64 // minimum distance between origin and destination, in meters
	usecase         *service.RideServiceMock
}

type startTripRequestResponse struct {
	Message string `json:"message"`
	Details struct {
		DecoderError []string `json:"decoder_error"`
	} `json:"details"`
}

func (s *startTripHandlerTestSuite) SetupTest() {
	s.usecase = &service.RideServiceMock{}
	s.minTripDistance = 5000

	defer mock.AssertExpectationsForObjects(s.T(), s.usecase)
}

func TestStartTripHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(startTripHandlerTestSuite))
}

// test form data
func (s *startTripHandlerTestSuite) TestEndpointValidatesFormDataCorrectly() {
	scenarios := []serializer.StartTripRequest{
		testutil.FakeStartTripRequestWithInvalidCoordinates(),
		testutil.FakeStartTripRequestWithInvalidUuid(),
	}

	for _, fakeReq := range scenarios {
		res := s.response(testutil.StartTripRequestToUrlValues(fakeReq))
		assert.Equal(s.T(), http.StatusBadRequest, res.Code)

		// Make sure the response is valid as well
		var result startTripRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			s.T().Errorf("failed to decode response body: %v", err)
		}

		// Make sure the response is for validation failed and not something else
		assert.Equal(s.T(), "validation failed", result.Message)
	}
}

// test if distance is less than minTripDistance, it should return false (and check response)
func (s *startTripHandlerTestSuite) TestMinimumTripDistanceIsMoreOrEqualToGivenMinDistance() {
	assert.False(s.T(), true)
}

// test if ride is in location
func (s *startTripHandlerTestSuite) TestRideIsInLocation() {
	assert.False(s.T(), true)
}

// test if ride.state is available or not
func (s *startTripHandlerTestSuite) TestRideIsAvailableBeforeStartingTrip() {
	assert.False(s.T(), true)
}

// test response
// 1. should have a specific format, should contain routes array
func (s *startTripHandlerTestSuite) TestTripStartsSuccessfullyGivenValidData() {
	assert.False(s.T(), true)
}

func (s *startTripHandlerTestSuite) response(data url.Values) *httptest.ResponseRecorder {
	rr := testutil.NewAPIResponseBuilder().
		Method(http.MethodPost).
		URL("/api/v1/trip/start").
		Body(data).
		Handler(NewUpdateRideLocationHandler(s.usecase).ServeHTTP).
		IsFormUrlEncoded().
		Build()

	return rr
}
