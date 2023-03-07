package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_healthcheckHandler(t *testing.T) {
	tests := []struct {
		name               string
		url                string
		method             string
		expectedStatusCode int
	}{
		{"health check success", "/api/v1/healthcheck", http.MethodGet, http.StatusOK},
	}
	var app application
	handler := http.HandlerFunc(app.healthCheckHandler)

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			req, err := http.NewRequest(e.method, e.url, nil)
			if err != nil {
				log.Fatal(err)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Result().StatusCode != e.expectedStatusCode {
				t.Errorf("%s: expected status %d, but got %d", e.name, e.expectedStatusCode, response.Result().StatusCode)
			}
		})

	}
}
