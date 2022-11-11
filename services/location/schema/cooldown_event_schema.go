package schema

import (
	"time"
)

type RideCooldownEvent struct {
	RideUuid  string        `json:"ride_uuid"`
	Timestamp int64         `json:"timestamp"`
	Duration  time.Duration `json:"-"`
}
