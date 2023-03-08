package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]any

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	MAX_BYTES := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(MAX_BYTES))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dest)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("the body should contain a single json object")
	}
	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, data envelope, statusCode int, headers http.Header) error {
	resBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for k, v := range headers {
		w.Header()[k] = v
	}
	resBytes = append(resBytes, '\n')
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resBytes)
	return nil
}

func (app *application) parseReqParam(r *http.Request, field string) (int64, error) {
	idStr := chi.URLParam(r, field)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
