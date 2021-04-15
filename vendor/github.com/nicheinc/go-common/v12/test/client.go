package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Server struct {
	*httptest.Server
	Hits int
}

type RequestFn func(r *http.Request, t *testing.T)

func TestServer(statusCode int, body interface{}, fn RequestFn, t *testing.T) *Server {
	s := &Server{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Hits++

		if fn != nil {
			fn(r, t)
		}

		w.WriteHeader(statusCode)
		switch body := body.(type) {
		case nil:
			break
		case []byte:
			w.Write(body)
		case string:
			w.Write([]byte(body))
		default:
			j, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("Error encoding body to JSON: %v", err)
			}
			w.Write(j)
		}
	}))
	return s
}

func (s *Server) ValidateHits(expectedHits int, t *testing.T) {
	if s.Hits != expectedHits {
		t.Errorf("Expected hits: %d, Actual: %d", expectedHits, s.Hits)
	}
}

func ValidateRequestMethod(expectedMethod string) RequestFn {
	return func(r *http.Request, t *testing.T) {
		ValidateMethod(r, expectedMethod, t)
	}
}

func ValidateRequestURL(expectedURL string) RequestFn {
	return func(r *http.Request, t *testing.T) {
		ValidateURL(r, expectedURL, t)
	}
}

func ValidateRequestHeaders(expectedHeaders http.Header) RequestFn {
	return func(r *http.Request, t *testing.T) {
		ValidateHeaders(r, expectedHeaders, t)
	}
}

func ValidateRequestBody(expectedBody interface{}) RequestFn {
	return func(r *http.Request, t *testing.T) {
		var expected string
		switch expectedBody := expectedBody.(type) {
		case nil:
			break
		case []byte:
			expected = string(expectedBody)
		case string:
			expected = expectedBody
		default:
			bytes, err := json.Marshal(expectedBody)
			if err != nil {
				t.Fatalf("Error JSON encoding expected request body: %s", err)
			}
			expected = string(bytes)
		}

		ValidateBody(r, expected, t)
	}
}

// Legacy methods for testing request fields directly

func ValidateMethod(r *http.Request, expectedMethod string, t *testing.T) {
	if r.Method != expectedMethod {
		t.Errorf("Expected method: '%s', Actual: '%s'", expectedMethod, r.Method)
	}
}

func ValidateURL(r *http.Request, expectedURL string, t *testing.T) {
	if r.RequestURI != expectedURL {
		t.Errorf("Expected URL: '%s', Actual: '%s'", expectedURL, r.RequestURI)
	}
}

func ValidateHeaders(r *http.Request, expectedHeaders http.Header, t *testing.T) {
	for key, expectedValues := range expectedHeaders {
		actualValues := r.Header[key]
		for _, expectedValue := range expectedValues {
			if !contains(actualValues, expectedValue) {
				t.Errorf("Expected header: '%s' to have value: '%s'", key, expectedValue)
			}
		}
	}
}

func ValidateBody(r *http.Request, expectedBody string, t *testing.T) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error reading request body: '%s'", err)
	}
	actualBody := string(bytes)

	if actualBody != expectedBody {
		t.Errorf("Expected body:\n%v\nActual:\n%v", expectedBody, actualBody)
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
