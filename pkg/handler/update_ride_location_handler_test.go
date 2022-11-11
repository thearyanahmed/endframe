package handler

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thearyanahmed/nordsec/pkg/presenter"
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

type updateRideLocationFailedValidationResponse struct {
	Message string `json:"message"`
	Details struct {
		Latitude  []string `json:"latitude"`
		Longitude []string `json:"longitude"`
		RideUuid  []string `json:"ride_uuid"`
	} `json:"details"`
}

func TestUpdateRideLocationHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateRideLocationHandlerTestSuite))
}

func (s *UpdateRideLocationHandlerTestSuite) SetupTest() {
	s.usecase = &service.RideServiceMock{}

	defer mock.AssertExpectationsForObjects(s.T(), s.usecase)
}

func (s *UpdateRideLocationHandlerTestSuite) TestFormRequestHandlesLatLonValidation() {
	data := testutil.FakeRecordRideEventRequestWithInvalidLatLon()
	res := s.response(testutil.RecordRideEventToUrlValues(data))

	var result updateRideLocationFailedValidationResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), http.StatusBadRequest, res.Code)
}

func (s *UpdateRideLocationHandlerTestSuite) TestFormRequestHandlesRideUuidValidation() {
	data := testutil.FakeRecordRideEventRequestWithInvalidRideUuid()
	res := s.response(testutil.RecordRideEventToUrlValues(data))

	var result updateRideLocationFailedValidationResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), http.StatusBadRequest, res.Code)
}

func (s *UpdateRideLocationHandlerTestSuite) TestRequestFailsWithoutFormData() {
	res := s.response(nil)

	assert.Equal(s.T(), http.StatusBadRequest, res.Code)

	var result updateRideLocationFailedValidationResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), "validation failed", result.Message)
	assert.Contains(s.T(), result.Details.Latitude, "The latitude field is required")
	assert.Contains(s.T(), result.Details.Longitude, "The longitude field is required")
	assert.Contains(s.T(), result.Details.RideUuid, "The ride_uuid field is required")
}

func (s *UpdateRideLocationHandlerTestSuite) TestRideLocationUpdatesSuccessfully() {
	data := testutil.FakeRecordRideEventRequest()

	s.usecase.On("UpdateRideLocation").Return(*(data.ToRideEvent()), nil).Once()

	res := s.response(testutil.RecordRideEventToUrlValues(data))

	assert.Equal(s.T(), http.StatusOK, res.Code)

	var result presenter.RideLocationUpdateResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		s.T().Errorf("failed to decode response body: %v", err)
	}

	assert.Equal(s.T(), "ride location updated", result.Message)
	assert.Equal(s.T(), result.Event.Uuid, data.RideUuid)
	assert.Equal(s.T(), result.Event.Latitude, fmt.Sprintf("%.6f", data.Latitude))
	assert.Equal(s.T(), result.Event.Longitude, fmt.Sprintf("%.6f", data.Longitude))
}

func (s *UpdateRideLocationHandlerTestSuite) response(data url.Values) *httptest.ResponseRecorder {
	rr := testutil.NewAPIResponseBuilder().
		Method(http.MethodPost).
		URL("/api/v1/ride/activate").
		Body(data).
		Handler(NewUpdateRideLocationHandler(s.usecase).ServeHTTP).
		IsFormUrlEncoded().
		Build()

	return rr
}
