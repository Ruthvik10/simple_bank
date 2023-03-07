package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]any{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.writeJSON(w, envelope{"error": message}, status, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundRespose(w http.ResponseWriter, r *http.Request) {
	messsage := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, messsage)
}
