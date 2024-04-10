package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ryanpujo/blog-app/internal/response"
	"github.com/stretchr/testify/require"
)

// HttpTest encapsulates the data needed to create an HTTP test case.
type HttpTest struct {
	URI        string
	BaseURI    string
	JSON       []byte
	HttpMethod string
}

// httpTestOption defines a function signature for options to modify an HttpTest instance.
type httpTestOption func(*HttpTest)

// WithBaseUri is an option setter for setting the base URI of an HttpTest instance.
func WithBaseUri(baseUri string) httpTestOption {
	return func(ht *HttpTest) {
		ht.BaseURI = baseUri
	}
}

// WithJson is an option setter for setting the JSON body of an HttpTest instance.
func WithJson(json []byte) httpTestOption {
	return func(ht *HttpTest) {
		ht.JSON = json
	}
}

// NewHttpTest creates a new HttpTest instance with the provided options.
func NewHttpTest(httpMethod, uri string, opts ...httpTestOption) *HttpTest {
	httpTest := &HttpTest{
		URI:        uri,
		HttpMethod: httpMethod,
	}

	// Apply each option to the HttpTest instance.
	for _, opt := range opts {
		opt(httpTest)
	}

	return httpTest
}

// ExecuteTest performs the HTTP test and returns the response and status code.
func (ht *HttpTest) ExecuteTest(t *testing.T, mux http.Handler) (*response.Response, int) {
	var body io.Reader

	// If a JSON body is provided, create a reader for it.
	if ht.JSON != nil {
		body = bytes.NewReader(ht.JSON)
	}

	// Construct the full request URI.
	fullURI := fmt.Sprintf("%s%s", ht.BaseURI, ht.URI)
	log.Println("Executing test for URI:", fullURI)

	// Create a new HTTP request with the provided method, URI, and body.
	req, err := http.NewRequest(ht.HttpMethod, fullURI, body)
	require.NoError(t, err)

	// Record the HTTP response using httptest.
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	// Decode the JSON response.
	var jsonRes *response.Response
	json.NewDecoder(recorder.Body).Decode(&jsonRes)
	// require.NoError(t, err)

	return jsonRes, recorder.Code
}
