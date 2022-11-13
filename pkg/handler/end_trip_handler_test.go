package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
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

func (h *endTripHandlerTestSuite) TestItFailsWithoutProperFormData() {
	scenarios := []serializer.EndTripRequest{
		testutil.FakeEndTripRequestWithInvalidUuid(),
		testutil.FakeEndTripRequestWithInvalidCoordinates(),
		testutil.FakeEndTripRequestWithMissingRequiredField(),
	}

	for _, scene := range scenarios {
		res := h.response(testutil.EndTripRequestToUrlValues(scene))

		assert.Equal(h.T(), http.StatusBadRequest, res.Code)
		var result testutil.Response
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			h.T().Errorf("failed to decode response body: %v", err)
		}

		// Make sure the response is for validation failed and not something else
		assert.Equal(h.T(), "validation failed", result.Message)
	}
}

func (h *endTripHandlerTestSuite) TestItFailsIfTripDoesNotExist() {
	scene := testutil.FakeEndTripRequest()

	h.rideService.On("GetRideEventByUuid").Return(entity.Event{}, errors.New("ride event not found")).Once()
	defer h.rideService.ResetMock()

	res := h.response(testutil.EndTripRequestToUrlValues(scene))

	assert.Equal(h.T(), http.StatusNotFound, res.Code)
}

func (h *endTripHandlerTestSuite) TestItFailsIfRideEventHasAlreadyEnded() {
	scene := testutil.FakeEndTripRequest()

	fakeEvent := testutil.FakeRideStatusRoaming()
	fakeEvent.TripUuid = scene.TripUuid
	fakeEvent.PassengerUuid = ""

	h.rideService.On("GetRideEventByUuid").Return(fakeEvent, nil).Once()
	h.rideService.On("TripHasEnded").Return(true).Once()
	defer h.rideService.ResetMock()

	res := h.response(testutil.EndTripRequestToUrlValues(scene))

	assert.Equal(h.T(), http.StatusBadRequest, res.Code)
}

func (h *endTripHandlerTestSuite) TestItFailsIfRideEventUuidDoesNotMatch() {
	scene := testutil.FakeEndTripRequest()

	fakeEvent := testutil.FakeRideStatusRoaming()

	differentFakeEvent := fakeEvent
	differentFakeEvent.TripUuid = uuid.New().String()

	h.rideService.On("GetRideEventByUuid").Return(differentFakeEvent, nil).Once()
	h.rideService.On("TripHasEnded").Return(false).Once()
	defer h.rideService.ResetMock()

	res := h.response(testutil.EndTripRequestToUrlValues(scene))

	assert.Equal(h.T(), http.StatusUnprocessableEntity, res.Code)

	var result testutil.Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		h.T().Errorf("failed to decode response body: %v", err)
	}

	// Make sure the response is for validation failed and not something else
	assert.Equal(h.T(), "given trip uuid does not match current trip uuid", result.Message)
}

func (h *endTripHandlerTestSuite) TestItEndsATripSuccessfully() {
	scene := testutil.FakeEndTripRequest()
	fakeEvent := testutil.FakeRideStatusRoaming()
	fakeEvent.TripUuid = scene.TripUuid

	h.rideService.On("GetRideEventByUuid").Return(fakeEvent, nil).Once()
	h.rideService.On("TripHasEnded").Return(false).Once()
	h.rideService.On("RecordEndRideEvent").Return(fakeEvent, nil).Once()
	h.rideService.On("EnterCooldownMode").Return(nil).Once()
	defer h.rideService.ResetMock()

	res := h.response(testutil.EndTripRequestToUrlValues(scene))

	assert.Equal(h.T(), http.StatusOK, res.Code)
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
