package testutil

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
	"net/url"
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

func RecordRideEventToUrlValues(rre serializer.RecordRideEventRequest) url.Values {
	req := url.Values{}

	req.Set("ride_uuid", rre.RideUuid)
	req.Set("latitude", fmt.Sprintf("%.6f", rre.Latitude))
	req.Set("longitude", fmt.Sprintf("%.6f", rre.Longitude))

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

func FakeStartTripRequestWithInvalidDistance() serializer.StartTripRequest {
	r := FakeStartTripRequest()

	// set a very minimum distance
	r.DestinationLatitude = r.OriginLatitude + 0.0006
	r.DestinationLongitude = r.OriginLongitude + 0.0006

	return r
}

func StartTripRequestToUrlValues(rre serializer.StartTripRequest) url.Values {
	req := url.Values{}

	req.Set("ride_uuid", rre.RideUuid)
	req.Set("client_uuid", rre.ClientUuid)
	req.Set("origin_latitude", fmt.Sprintf("%.6f", rre.OriginLatitude))
	req.Set("origin_longitude", fmt.Sprintf("%.6f", rre.OriginLongitude))
	req.Set("destination_latitude", fmt.Sprintf("%.6f", rre.DestinationLatitude))
	req.Set("destination_longitude", fmt.Sprintf("%.6f", rre.DestinationLongitude))

	return req
}
