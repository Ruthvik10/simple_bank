package store

import (
	"database/sql"
	"errors"

	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/jmoiron/sqlx"
)

type AccountStore struct {
	db *sqlx.DB
}

func (store AccountStore) Get(id int64) (*models.Account, error) {
	var acc models.Account
	err := store.db.Get(&acc, "SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &acc, nil
}

func (store AccountStore) Create(acc *models.Account) error {
	query := `INSERT INTO accounts (owner, balance, currency) VALUES ($1, $2, $3) RETURNING *`
	err := store.db.QueryRowx(query, acc.Owner, acc.Balance, acc.Currency).StructScan(acc)
	if err != nil {
		return err
	}
	return nil
}

func (store AccountStore) UpdateAccount(acc *models.Account) error {
	query := `UPDATE accounts SET owner=$1, balance=$2, currency=$3 WHERE id=$4 RETURNING *`
	args := []any{acc.Owner, acc.Balance, acc.Currency, acc.ID}
	err := store.db.QueryRowx(query, args...).StructScan(acc)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (store AccountStore) UpdateBalance(acc *models.Account) error {
	query := `UPDATE accounts SET balance=$1 WHERE id=$2 RETURNING *`
	err := store.db.QueryRowx(query, acc.Balance, acc.ID).StructScan(acc)
	if err != nil {
		return err
	}
	return nil
}

func (store AccountStore) List() ([]*models.Account, error) {
	var accounts []*models.Account
	query := `SELECT * FROM accounts ORDER BY id ASC`
	err := store.db.Select(&accounts, query)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (store AccountStore) Delete(id int64) error {
	query := `DELETE FROM accounts WHERE id=$1`
	result, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if nRows == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (store AccountStore) GetForUpdate(id int64) (*models.Account, error) {
	var acc models.Account
	err := store.db.Get(&acc, "SELECT * FROM accounts WHERE id = $1 FOR NO KEY UPDATE", id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &acc, nil
}
