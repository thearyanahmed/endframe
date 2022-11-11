package testutil

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
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
