package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/pkg/service"
	"github.com/thearyanahmed/nordsec/pkg/testutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type notifyPositionHandlerTestSuite struct {
	suite.Suite
	rideService *service.RideServiceMock
}

type notifyTripLocationSuccessResponse struct {
	Uuid          string  `json:"uuid"`
	RideUuid      string  `json:"ride_uuid"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	PassengerUuid string  `json:"passenger_uuid"`
	TripUuid      string  `json:"trip_uuid"`
	Timestamp     int64   `json:"timestamp"`
	State         string  `json:"state"`
}

type notifyTripLocationFailedValidationResponse struct {
	Message string `json:"message"`
	Details struct {
		DecoderError []string `json:"decoder_error"`
	} `json:"details"`
}

func (s *notifyPositionHandlerTestSuite) SetupTest() {
	s.rideService = &service.RideServiceMock{}

	defer mock.AssertExpectationsForObjects(s.T(), s.rideService)
}

func TestNotifyPositionHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(notifyPositionHandlerTestSuite))
}

func (s *notifyPositionHandlerTestSuite) TestEndpointValidatesFormDataCorrectly() {
	scenarios := []serializer.NotifyTripLocationRequest{
		testutil.FakeNotifyTripLocationRequestWithInvalidLatLon(),
		testutil.FakeNotifyTripLocationRequestWithMissingRequiredFields(),
		testutil.FakeNotifyTripLocationRequestWithInvalidUuid(),
	}

	for _, fakeReq := range scenarios {
		res := s.response(testutil.NotifyTripLocationRequestToUrlValues(fakeReq))
		assert.Equal(s.T(), http.StatusBadRequest, res.Code)

		// Make sure the response is valid as well
		var result notifyTripLocationFailedValidationResponse
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			s.T().Errorf("failed to decode response body: %v", err)
		}

		// Make sure the response is for validation failed and not something else
		assert.Equal(s.T(), "validation failed", result.Message)
	}
}

func (s *notifyPositionHandlerTestSuite) TestRideLocationUpdatesSuccessfullyWithValidData() {
	fakeReq := testutil.FakeNotifyTripLocationRequest()
	formData := testutil.NotifyTripLocationRequestToUrlValues(fakeReq)
	rideEvent := fakeReq.ToRideEvent()

	s.rideService.On("UpdateRideLocation").Return(rideEvent, nil).Once()
	res := s.response(formData)
	assert.Equal(s.T(), http.StatusOK, res.Code)

	var result notifyTripLocationSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	// Make sure the response is for validation failed and not something else
	assert.Equal(s.T(), rideEvent.TripUuid, result.TripUuid)
	assert.Equal(s.T(), rideEvent.Lat, result.Lat)
	assert.Equal(s.T(), rideEvent.Lon, result.Lon)
	assert.Equal(s.T(), rideEvent.PassengerUuid, result.PassengerUuid)
	assert.Equal(s.T(), rideEvent.RideUuid, result.RideUuid)
	assert.NotEmpty(s.T(), result.Timestamp)
}

func (s *notifyPositionHandlerTestSuite) response(data url.Values) *httptest.ResponseRecorder {
	rr := testutil.NewAPIResponseBuilder().
		Method(http.MethodPost).
		URL("/api/v1/trip/notify/location").
		Body(data).
		Handler(NewNotifyPositionHandler(s.rideService).ServeHTTP).
		IsFormUrlEncoded().
		Build()

	return rr
}
