package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {
	var app application
	var registered = []struct {
		route  string
		method string
	}{
		{"/api/v1/healthcheck", "GET"},
		{"/api/v1/accounts/", "POST"},
		{"/api/v1/accounts/{id:^[0-9]+}", "GET"},
		{"/api/v1/accounts/", "PUT"},
		{"/api/v1/accounts/{id:^[0-9]+}", "PATCH"},
	}
	routes := app.routes()

	chiRoutes := routes.(chi.Routes)
	for _, route := range registered {
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("error: route %s, method %s is not registered", route.route, route.method)
		}
	}
}

func routeExists(testRoute string, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if testRoute == route && testMethod == method {
			found = true
		}
		return nil
	})
	return found
}
