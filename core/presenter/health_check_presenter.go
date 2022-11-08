package presenter

import "time"

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func HealthCheck(status string, t time.Time) *HealthCheckResponse {
	return &HealthCheckResponse{
		Status:    status,
		Timestamp: t,
	}
}
