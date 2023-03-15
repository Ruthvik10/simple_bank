package main

import (
	"errors"
	"net/http"

	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/Ruthvik10/simple_bank/internal/store"
)

func (app *application) createTransferHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FromAccountID int64 `json:"from_account_id"`
		ToAccountID   int64 `json:"to_account_id"`
		Amount        int64 `json:"amount"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	transfer := &models.Transfer{
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		Amount:        input.Amount,
	}

	err = app.store.Transfer.CreateTransfer(transfer)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrInvalidPayer), errors.Is(err, store.ErrInvalidPayee), errors.Is(err, store.ErrInsufficientBalance):
			app.badRequestErrorResponse(w, r, err)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	err = app.writeJSON(w, envelope{"message": "successfully transferred the ammount"}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
