package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/thearyanahmed/endframe/pkg/service/ride"
	"github.com/thearyanahmed/endframe/pkg/testutil"
)

type nearByRidesHandlerTestSuite struct {
	suite.Suite
	rideService *ride.RideServiceMock
}

type nearByRidesSuccessResponse struct {
	RideUuid string  `json:"ride_uuid"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	State    string  `json:"state"`
}

func (s *nearByRidesHandlerTestSuite) SetupTest() {
	s.rideService = &ride.RideServiceMock{}

	defer mock.AssertExpectationsForObjects(s.T(), s.rideService)
}

func TestNearByRidesHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(nearByRidesHandlerTestSuite))
}

func (s *nearByRidesHandlerTestSuite) TestValidateCoordinateParamsAreRequired() {
	scenarios := []string{
		"",
		"x1=5000&y1=50000&x3=10000&y3=100000",
	}

	for _, scenario := range scenarios {
		res := s.response(scenario) // no query params
		assert.Equal(s.T(), http.StatusBadRequest, res.Code)

		var result testutil.Response
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			s.T().Errorf("failed to decode response body: %v", err)
		}

		// Make sure the response is for validation failed and not something else
		assert.Equal(s.T(), "validation failed", result.Message)
	}
}

func (s *nearByRidesHandlerTestSuite) TestItReturnsAListOfRidesWhenGivenValidData() {
	ridesCount := 10

	rides := testutil.FakeRideEntity(ridesCount, "")
	s.rideService.On("FindNearByRides").Return(rides, nil).Once()
	res := s.response("x1=52.3251&y1=13.453&x3=52.3361&y3=13.475")

	assert.Equal(s.T(), http.StatusOK, res.Code)

	var result []nearByRidesSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), ridesCount, len(result))
}

func (s *nearByRidesHandlerTestSuite) TestRidesCanBeFilteredWithAValidState() {
	ridesCount := 10
	state := "in_route"

	rides := testutil.FakeRideEntity(ridesCount, state)
	s.rideService.On("FindNearByRides").Return(rides, nil).Once()
	res := s.response(fmt.Sprintf("x1=52.3251&y1=13.453&x3=52.3361&y3=13.475&state=%s", state))

	assert.Equal(s.T(), http.StatusOK, res.Code)

	var result []nearByRidesSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), ridesCount, len(result))

	for _, r := range result {
		assert.Equal(s.T(), state, r.State)
	}
}

func (s *nearByRidesHandlerTestSuite) TestProvidingInvalidStateInFilterReturnsAllRidesInThatLocation() {
	ridesCount := 10
	state := "some-invalid-state"

	rides := testutil.FakeRideEntity(ridesCount, "")
	s.rideService.On("FindNearByRides").Return(rides, nil).Once()
	res := s.response(fmt.Sprintf("x1=52.3251&y1=13.453&x3=52.3361&y3=13.475&state=%s", state))

	assert.Equal(s.T(), http.StatusOK, res.Code)

	var result []nearByRidesSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), ridesCount, len(result))
}

func (s *nearByRidesHandlerTestSuite) response(query string) *httptest.ResponseRecorder {
	rr := testutil.NewAPIResponseBuilder().
		Method(http.MethodGet).
		URL(fmt.Sprintf("/api/v1/rides/near-by?%s", query)).
		Handler(NewNearByRidesHandler(s.rideService).ServeHTTP).
		Build()

	return rr
}
