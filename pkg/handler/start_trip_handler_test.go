package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/services/location/entity"
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
	rideService     *service.RideServiceMock
}

type startTripRequestValidationResponse struct {
	Message string `json:"message"`
	Details struct {
		DecoderError []string `json:"decoder_error"`
	} `json:"details"`
}

func (s *startTripHandlerTestSuite) SetupTest() {
	s.rideService = &service.RideServiceMock{}
	s.minTripDistance = 5000

	defer mock.AssertExpectationsForObjects(s.T(), s.rideService)
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
		var result startTripRequestValidationResponse
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			s.T().Errorf("failed to decode response body: %v", err)
		}

		// Make sure the response is for validation failed and not something else
		assert.Equal(s.T(), "validation failed", result.Message)
	}
}

// test if distance is less than minTripDistance, it should return false (and check response)
func (s *startTripHandlerTestSuite) TestTripDistanceIsMoreOrEqualToGivenMinDistance() {
	fakeReq := testutil.FakeStartTripRequest()

	minDist := 50000.00

	s.rideService.On("GetMinimumTripDistance").Return(minDist).Once()
	s.rideService.On("DistanceIsGreaterThanMinimumDistance").Return(false).Once()

	res := s.response(testutil.StartTripRequestToUrlValues(fakeReq))

	assert.Equal(s.T(), http.StatusBadRequest, res.Code)

	// Make sure the response is valid as well
	var result testutil.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	// Make sure the response is for validation failed and not something else
	msg := fmt.Sprintf("distance too low. minimum distance required %.2f meters or greater", minDist)
	assert.Equal(s.T(), msg, result.Message)
}

// test if ride is in location
func (s *startTripHandlerTestSuite) TestItReturnsErrorIfRideIsNotInLocation() {
	fakeReq := testutil.FakeStartTripRequest()

	riderErr := errors.New("ride not in nearby location")

	s.rideService.On("GetMinimumTripDistance").Return(1).Once() // a very low distance
	s.rideService.On("DistanceIsGreaterThanMinimumDistance").Return(true).Once()
	s.rideService.On("FindRideInLocation").Return(entity.Ride{}, riderErr).Once()

	res := s.response(testutil.StartTripRequestToUrlValues(fakeReq))
	assert.Equal(s.T(), http.StatusUnprocessableEntity, res.Code)

	var result testutil.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), riderErr.Error(), result.Message)
}

// test if ride.state is available or not
func (s *startTripHandlerTestSuite) TestRideIsAvailableBeforeStartingTrip() {
}

// test response
// 1. should have a specific format, should contain routes array
func (s *startTripHandlerTestSuite) TestTripStartsSuccessfullyGivenValidData() {
}

func (s *startTripHandlerTestSuite) response(data url.Values) *httptest.ResponseRecorder {
	rr := testutil.NewAPIResponseBuilder().
		Method(http.MethodPost).
		URL("/api/v1/trip/start").
		Body(data).
		Handler(NewStartTripHandler(s.rideService).ServeHTTP).
		IsFormUrlEncoded().
		Build()

	return rr
}
