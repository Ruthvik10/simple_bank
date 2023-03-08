package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthCheckHandler)
		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", app.CreateAccountHandler)
			r.Get("/{id:^[0-9]+}", app.getAccountByIDHandler)
			r.Put("/", app.updateAccountHandler)
			r.Patch("/{id:^[0-9]+}", app.updateBalanceHandler)
			r.Get("/", app.listAccountsHandler)
			r.Delete("/{id:^[0-9]+}", app.deleteAccountHandler)
		})
	})
	return r
}
