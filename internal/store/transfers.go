package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/jmoiron/sqlx"
)

type TransferStore struct {
	db *sqlx.DB
}

var (
	ErrInsufficientBalance = errors.New("insufficent balance")
	ErrInvalidPayer        = errors.New("invalid payer details")
	ErrInvalidPayee        = errors.New("invalid payee details")
)

func (store TransferStore) CreateTransfer(t *models.Transfer) error {
	tx, err := store.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// read from_account info
	var fromAccount models.Account
	err = tx.Get(&fromAccount, "SELECT * FROM accounts WHERE id=$1 FOR NO KEY UPDATE", t.FromAccountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrInvalidPayer
		default:
			return err
		}
	}
	if fromAccount.Balance < t.Amount {
		return ErrInsufficientBalance
	}

	// check if the to_account exists
	var toAccount models.Account
	err = tx.Get(&toAccount, "SELECT * FROM accounts WHERE id=$1 FOR NO KEY UPDATE", t.ToAccountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrInvalidPayee
		default:
			return err
		}
	}

	// create a transfer record
	_, err = tx.Exec("INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)", t.FromAccountID, t.ToAccountID, t.Amount)
	if err != nil {
		return nil
	}

	// create an entry for from_account to_account
	stmt, err := tx.Prepare(
		`INSERT INTO entries (account_id, amount) VALUES ($1, $2)`,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(t.FromAccountID, -t.Amount)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t.ToAccountID, t.Amount)
	if err != nil {
		return err
	}

	// update balance for from_account
	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id=$2", t.Amount, t.FromAccountID)
	if err != nil {
		return err
	}

	// update balance for to_account
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id=$2", t.Amount, t.ToAccountID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
