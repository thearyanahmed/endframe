package testutil

import (
	"net/url"
)

type Response struct {
	HttpStatusCode int        `json:"-"`
	Message        string     `json:"message"`
	Details        url.Values `json:"details,omitempty"`
}
