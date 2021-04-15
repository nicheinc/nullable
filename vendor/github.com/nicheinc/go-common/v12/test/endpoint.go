package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestRequest sends a mock request with the specified HTTP method, URL, headers,
// and request body to a test server initialized with the provided router. It will then
// assert that the response has the expected status code and response body. The expectedBody
// will be serialized to JSON before being compared with the actual response's body.
// NOTE: If nil is provided for the requestBody, the request will have no body. If
//		a string is provided, the body will be the string provided. If an object is
//		provided, it will be serialized to JSON before being set as the request body.
func TestRequest(method, url string, headers http.Header, requestBody interface{},
	expectedStatusCode int, expectedResponse interface{}, router http.Handler, t *testing.T) {

	res := MakeRequest(method, url, headers, requestBody, router, t)
	CheckResponse(res, expectedStatusCode, expectedResponse, t)
}

// TestRequestWithHeader sends a mock request with the specified HTTP method,
// URL, headers, and request body to a test server initialized with the
// provided router. It will then assert that the response has the expected
// headers (although it can other headers, as well), status code, and response
// body. The expectedBody will be serialized to JSON before being compared with
// the actual response's body.
// NOTE: If nil is provided for the requestBody, the request will have no body. If
//		a string is provided, the body will be the string provided. If an object is
//		provided, it will be serialized to JSON before being set as the request body.
func TestRequestWithHeader(method, url string, headers http.Header, requestBody interface{},
	expectedHeader http.Header, expectedStatusCode int, expectedResponse interface{},
	router http.Handler, t *testing.T) {

	res := MakeRequest(method, url, headers, requestBody, router, t)
	CheckResponseWithHeader(res, expectedHeader, expectedStatusCode, expectedResponse, t)
}

// MakeRequest sends a mock request with the specified HTTP method, URL, content-type,
// and request body to a test server initialized with the provided router. The
// response is returned.
// NOTE: If nil is provided for the requestBody, the request will have no body. If an
//		io.Reader, []byte, or string is provided, the body will use that. If anything else
//		is provided, it will be serialized to JSON before being set as the request body.
func MakeRequest(method, url string, headers http.Header, requestBody interface{},
	router http.Handler, t *testing.T) *http.Response {

	server := httptest.NewServer(router)
	defer server.Close()

	var bodyReader io.Reader
	switch requestBody := requestBody.(type) {
	case nil:
		break
	case io.Reader:
		bodyReader = requestBody
	case []byte:
		bodyReader = bytes.NewReader(requestBody)
	case string:
		bodyReader = bytes.NewReader([]byte(requestBody))
	default:
		actualBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Error JSON encoding request body")
		}
		bodyReader = bytes.NewReader(actualBytes)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", server.URL, url), bodyReader)
	if err != nil {
		t.Fatalf("Error creating new request: %s", err)
	}

	if headers != nil {
		req.Header = headers
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request to test server: %s", err)
	}

	return res
}

// CheckResponse asserts that the given response has the expected status code and
// body. The expectedBody that will be serialized to JSON before being compared with
// the actual response's body.
func CheckResponse(res *http.Response, expectedStatusCode int, expectedBody interface{}, t *testing.T) {

	// Get expected response:
	var expectedBytes []byte
	switch expectedBody := expectedBody.(type) {
	case nil:
		break
	case []byte:
		expectedBytes = expectedBody
	case string:
		expectedBytes = []byte(expectedBody)
	default:
		var err error
		expectedBytes, err = json.Marshal(expectedBody)
		if err != nil {
			t.Fatalf("Error marshalling expected response to JSON: %s", err)
		}
	}
	expectedResponse := string(expectedBytes)

	// Try to format/indent the expected response (if this fails, it might
	// just mean that the expected response isn't JSON, which is okay):
	var expectedBuf bytes.Buffer
	if err := json.Indent(&expectedBuf, expectedBytes, "", "    "); err == nil {
		expectedResponse = expectedBuf.String()
	}

	// Get the actual response:
	actualBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err)
	}
	actualResponse := string(actualBytes)

	// Try to format/indent the actual response (if this fails, it might
	// just mean that the expected response isn't JSON, which is okay):
	var actualBuf bytes.Buffer
	if err := json.Indent(&actualBuf, actualBytes, "", "    "); err == nil {
		actualResponse = actualBuf.String()
	}

	// Check the status code:
	if res.StatusCode != expectedStatusCode {
		t.Errorf("\nExpected Status Code: %d\nActual: %d", expectedStatusCode, res.StatusCode)
	}

	// Check the response body:
	if actualResponse != expectedResponse {
		// It's possible this is a false positive if the two are JSON objects with keys in
		// a different order. Attempt to unmarshal both and compare deeply.
		var actualVal interface{}
		var expectedVal interface{}
		if err := json.Unmarshal([]byte(actualResponse), &actualVal); err != nil {
			t.Errorf("\nExpected response:\n%s\nActual:\n%s", expectedResponse, actualResponse)
		} else if err := json.Unmarshal([]byte(expectedResponse), &expectedVal); err != nil {
			t.Errorf("\nExpected response:\n%s\nActual:\n%s", expectedResponse, actualResponse)
		} else if !reflect.DeepEqual(actualVal, expectedVal) {
			t.Errorf("\nExpected response:\n%s\nActual:\n%s", expectedResponse, actualResponse)
		}
	}
}

// CheckResponseWithHeader asserts that the given response has the expected
// headers, status code, and body. Note that the response may have other
// headers, besides the ones whose presence was verified. The expectedBody that
// will be serialized to JSON before being compared with the actual response's
// body.
func CheckResponseWithHeader(res *http.Response, expectedHeader http.Header, expectedStatusCode int, expectedBody interface{}, t *testing.T) {
	// Check response header
	for key, val := range expectedHeader {
		if !reflect.DeepEqual(res.Header[key], val) {
			t.Errorf("\nExpected %s Header: %v\nActual: %v", key, val, res.Header[key])
		}
	}

	CheckResponse(res, expectedStatusCode, expectedBody, t)
}
