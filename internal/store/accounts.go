package store

import (
	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/jmoiron/sqlx"
)

type AccountStore struct {
	DB *sqlx.DB
}

func (store AccountStore) Get(id int64) (*models.Account, error) {
	var acc models.Account
	err := store.DB.Get(&acc, "SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (store AccountStore) Create(acc *models.Account) error {
	query := `INSERT INTO accounts (owner, balance, currency) VALUES $1, $2, $3 RETURNING *`
	err := store.DB.QueryRowx(query, acc.Owner, acc.Balance, acc.Currency).StructScan(acc)
	if err != nil {
		return err
	}
	return nil
}
