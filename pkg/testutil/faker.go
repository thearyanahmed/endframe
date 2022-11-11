package testutil

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/thearyanahmed/nordsec/pkg/serializer"
)

func FakeRideEventRequest() serializer.RecordRideEventRequest {
	return serializer.RecordRideEventRequest{
		RideUuid:  gofakeit.UUID(),
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}