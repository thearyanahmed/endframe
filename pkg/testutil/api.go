package testutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type APIResponseBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    url.Values
	handler http.HandlerFunc
}

func NewAPIResponseBuilder() *APIResponseBuilder {
	return &APIResponseBuilder{
		headers: map[string]string{},
	}
}

func (b *APIResponseBuilder) Method(method string) *APIResponseBuilder {
	b.method = method
	return b
}

func (b *APIResponseBuilder) URL(url string) *APIResponseBuilder {
	b.url = url
	return b
}

func (b *APIResponseBuilder) Body(body url.Values) *APIResponseBuilder {
	b.body = body
	return b
}

func (b *APIResponseBuilder) Handler(handler http.HandlerFunc) *APIResponseBuilder {
	b.handler = handler
	return b
}

func (b *APIResponseBuilder) IsFormUrlEncoded() *APIResponseBuilder {
	b.headers["Content-Type"] = "application/x-www-form-urlencoded"
	return b
}

func (b *APIResponseBuilder) Build() *httptest.ResponseRecorder {
	req := httptest.NewRequest(b.method, b.url, strings.NewReader(b.body.Encode()))

	for k, v := range b.headers {
		req.Header.Add(k, v)
	}

	rr := httptest.NewRecorder()

	b.handler.ServeHTTP(rr, req)

	return rr
}
