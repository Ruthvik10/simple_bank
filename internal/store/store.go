package store

import (
	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	Account models.AccountStore
}

func NewStore(db *sqlx.DB) Store {
	return Store{
		Account: AccountStore{
			DB: db,
		},
	}
}
