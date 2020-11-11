package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setUp() {
	jsonPath := "credentials/test.json"
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", jsonPath)
	os.Setenv(envKeyLoggerName, "test-my-app")
	os.Setenv(envKeyLoggerJob, "api")
}

func TestApi(t *testing.T) {
	setUp()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

}

func TestGetOrCreateRequestIdWithoutRequestID(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	requestID := getOrCreateRequestID(req.Header)
	fmt.Print(requestID)
	if requestID == "" {
		t.Error("requestId should be set")
	}
}

func TestGetOrCreateRequestIdWithRequestID(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	expectedRequestID := "request_id"
	req.Header.Set(headerNameRequestID, expectedRequestID)
	requestID := getOrCreateRequestID(req.Header)
	fmt.Print(requestID)
	if requestID != expectedRequestID {
		t.Error("requestId should be set")
	}
}
