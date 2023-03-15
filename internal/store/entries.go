package store

import (
	"database/sql"
	"errors"

	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/jmoiron/sqlx"
)

type EntryStore struct {
	db *sqlx.DB
}

func (store EntryStore) Create(e *models.Entry) error {
	query := "INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING *"
	err := store.db.QueryRowx(query, e.Amount, e.AccountID).StructScan(e)
	if err != nil {
		return err
	}
	return nil
}

func (store EntryStore) Get(accountID int64) (*models.Entry, error) {
	var entry models.Entry
	err := store.db.Get(&entry, "SELECT * FROM entries WHERE account_id=$1", accountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &entry, nil
}
