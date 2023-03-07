package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string, 0)
	response["env"] = app.cfg.env
	response["status"] = "healthy"
	resBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": true, "message": "something went wrong, try again later"}`))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}
