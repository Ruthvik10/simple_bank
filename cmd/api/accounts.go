package main

import (
	"errors"
	"net/http"

	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/Ruthvik10/simple_bank/internal/store"
)

func (app *application) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Owner    string `json:"owner"`
		Balance  int64  `json:"balance"`
		Currency string `json:"currency"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}
	acc := &models.Account{
		Owner:    input.Owner,
		Balance:  input.Balance,
		Currency: input.Currency,
	}
	err = app.store.Account.Create(acc)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJSON(w, envelope{"account": acc}, http.StatusCreated, nil)
}

func (app *application) getAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseReqParam(r, "id")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	acc, err := app.store.Account.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundRespose(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	app.writeJSON(w, envelope{"account": acc}, http.StatusOK, nil)
}
