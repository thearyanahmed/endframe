package testutil

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"github.com/thearyanahmed/nordsec/pkg/service/location/entity"
	"net/url"
	"time"
)

func FakeRecordRideEventRequest() serializer.RecordRideEventRequest {
	return serializer.RecordRideEventRequest{
		RideUuid:  gofakeit.UUID(),
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}

func FakeRecordRideEventRequestWithInvalidLatLon() serializer.RecordRideEventRequest {
	return serializer.RecordRideEventRequest{
		RideUuid:  gofakeit.UUID(),
		Latitude:  523.012,
		Longitude: -320.22,
	}
}
func FakeRecordRideEventRequestWithInvalidRideUuid() serializer.RecordRideEventRequest {
	return serializer.RecordRideEventRequest{
		RideUuid:  "not-an-uuid",
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}

func RecordRideEventToUrlValues(request serializer.RecordRideEventRequest) url.Values {
	req := url.Values{}

	req.Set("ride_uuid", request.RideUuid)
	req.Set("latitude", fmt.Sprintf("%.6f", request.Latitude))
	req.Set("longitude", fmt.Sprintf("%.6f", request.Longitude))

	return req
}

func FakeStartTripRequest() serializer.StartTripRequest {
	originLat := gofakeit.Latitude()
	originLon := gofakeit.Longitude()

	destLat, destLon := originLat+1.0232, originLon+1.0232

	return serializer.StartTripRequest{
		RideUuid:             uuid.New().String(),
		ClientUuid:           uuid.New().String(),
		OriginLatitude:       originLat,
		OriginLongitude:      originLon,
		DestinationLatitude:  destLat,
		DestinationLongitude: destLon,
	}
}

func FakeStartTripRequestWithInvalidCoordinates() serializer.StartTripRequest {
	r := FakeStartTripRequest()

	r.DestinationLongitude = 325.12
	r.DestinationLatitude = 325.12

	r.OriginLatitude = 325.12
	r.OriginLongitude = 325.12

	return r
}

func FakeStartTripRequestWithInvalidUuid() serializer.StartTripRequest {
	r := FakeStartTripRequest()

	r.ClientUuid = "not-a-valid-uuid"
	r.RideUuid = "not-a-valid-uuid"

	return r
}

func StartTripRequestToUrlValues(request serializer.StartTripRequest) url.Values {
	req := url.Values{}

	req.Set("ride_uuid", request.RideUuid)
	req.Set("client_uuid", request.ClientUuid)
	req.Set("origin_latitude", fmt.Sprintf("%.6f", request.OriginLatitude))
	req.Set("origin_longitude", fmt.Sprintf("%.6f", request.OriginLongitude))
	req.Set("destination_latitude", fmt.Sprintf("%.6f", request.DestinationLatitude))
	req.Set("destination_longitude", fmt.Sprintf("%.6f", request.DestinationLongitude))

	return req
}

func FakeRoute(origin, dest entity.Coordinate) []entity.Coordinate {
	var route []entity.Coordinate

	route = append(route, origin)

	for i := 1; i < 11; i++ {
		point := entity.Coordinate{
			Lat: gofakeit.Latitude() + 0.005,
			Lon: gofakeit.Longitude() + 0.005,
		}
		route = append(route, point)
	}

	route = append(route, dest)

	return route
}

func FakeRoamingRideEntity() entity.Ride {
	return entity.Ride{
		RideUuid: uuid.New().String(),
		Lat:      gofakeit.Latitude(),
		Lon:      gofakeit.Longitude(),
		State:    entity.StateRoaming,
	}
}

func FakeInRouteRideEntity() entity.Ride {
	return entity.Ride{
		RideUuid: uuid.New().String(),
		Lat:      gofakeit.Latitude(),
		Lon:      gofakeit.Longitude(),
		State:    entity.StateInRoute,
	}
}

func FakeEventInRoute(rideUuid string, loc entity.Coordinate) entity.Event {
	return entity.Event{
		Uuid:          uuid.New().String(),
		RideUuid:      rideUuid,
		Lat:           loc.Lat,
		Lon:           loc.Lon,
		PassengerUuid: uuid.New().String(),
		TripUuid:      uuid.New().String(),
		Timestamp:     gofakeit.Int64(),
		State:         entity.StateInRoute,
	}
}

func FakeNotifyTripLocationRequest() serializer.NotifyTripLocationRequest {
	return serializer.NotifyTripLocationRequest{
		Latitude:      gofakeit.Latitude(),
		Longitude:     gofakeit.Longitude(),
		RideUuid:      uuid.New().String(),
		TripUuid:      uuid.New().String(),
		ClientUuid:    uuid.New().String(),
		PassengerUuid: uuid.New().String(),
	}
}

func FakeNotifyTripLocationRequestWithInvalidLatLon() serializer.NotifyTripLocationRequest {
	r := FakeNotifyTripLocationRequest()
	r.Latitude = 122012.22
	r.Longitude = 122012.22

	return r
}

func FakeNotifyTripLocationRequestWithMissingRequiredFields() serializer.NotifyTripLocationRequest {
	return serializer.NotifyTripLocationRequest{}
}

func FakeNotifyTripLocationRequestWithInvalidUuid() serializer.NotifyTripLocationRequest {
	r := FakeNotifyTripLocationRequest()

	r.PassengerUuid = "not-a-valid-uuid"
	r.ClientUuid = "not-a-valid-uuid"
	r.TripUuid = "not-a-valid-uuid"

	return r
}

func NotifyTripLocationRequestToUrlValues(request serializer.NotifyTripLocationRequest) url.Values {
	req := url.Values{}

	req.Set("ride_uuid", request.RideUuid)
	req.Set("trip_uuid", request.TripUuid)
	req.Set("client_uuid", request.ClientUuid)
	req.Set("passenger_uuid", request.PassengerUuid)
	req.Set("latitude", fmt.Sprintf("%.6f", request.Latitude))
	req.Set("longitude", fmt.Sprintf("%.6f", request.Longitude))

	return req
}

func Lat() float64 {
	return gofakeit.Latitude()
}

func Lon() float64 {
	return gofakeit.Longitude()
}

func FakeRideEntity(n int, targetState string) []entity.Ride {
	var rides []entity.Ride

	for i := 0; i < n; i++ {
		state := targetState
		if targetState == "" {
			state = FakeRandomState()
		}

		rides = append(rides, entity.Ride{
			RideUuid: uuid.New().String(),
			Lat:      Lat(),
			Lon:      Lon(),
			State:    state,
		})
	}

	return rides
}

func FakeRandomState() string {
	switch gofakeit.IntRange(0, 2) {
	case 0:
		return entity.StateRoaming
	case 1:
		return entity.StateInRoute
	case 2:
		return entity.StateInCooldown
	}

	return entity.StateRoaming
}

func FakeEndTripRequest() serializer.EndTripRequest {
	return serializer.EndTripRequest{
		RideUuid:      uuid.New().String(),
		Latitude:      gofakeit.Latitude(),
		Longitude:     gofakeit.Longitude(),
		ClientUuid:    uuid.New().String(),
		PassengerUuid: uuid.New().String(),
		TripUuid:      uuid.New().String(),
	}
}

func FakeEndTripRequestWithInvalidUuid() serializer.EndTripRequest {
	r := FakeEndTripRequest()

	r.ClientUuid = "not-an-uuid"
	r.RideUuid = "not-an-uuid"
	r.PassengerUuid = "not-an-uuid"

	return r
}

func FakeEndTripRequestWithInvalidCoordinates() serializer.EndTripRequest {
	r := FakeEndTripRequest()

	r.Latitude = 1203123.123
	r.Longitude = 1203123.123

	return r
}

func FakeEndTripRequestWithMissingRequiredField() serializer.EndTripRequest {
	r := FakeEndTripRequest()

	r.ClientUuid = ""
	r.PassengerUuid = ""
	r.RideUuid = ""

	return r
}

func EndTripRequestToUrlValues(request serializer.EndTripRequest) url.Values {
	req := url.Values{}

	req.Set("ride_uuid", request.RideUuid)
	req.Set("trip_uuid", request.TripUuid)
	req.Set("client_uuid", request.ClientUuid)
	req.Set("passenger_uuid", request.PassengerUuid)
	req.Set("latitude", fmt.Sprintf("%.6f", request.Latitude))
	req.Set("longitude", fmt.Sprintf("%.6f", request.Longitude))

	return req
}

func FakeRideStatusRoaming() entity.Event {
	return entity.Event{
		Uuid:          uuid.New().String(),
		RideUuid:      uuid.New().String(),
		Lat:           gofakeit.Latitude(),
		Lon:           gofakeit.Longitude(),
		PassengerUuid: uuid.New().String(),
		TripUuid:      uuid.New().String(),
		Timestamp:     time.Now().Unix(),
		State:         entity.StateRoaming,
	}
}
