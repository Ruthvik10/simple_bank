package store

import (
	"errors"

	mock "github.com/Ruthvik10/simple_bank/internal/mock/db"
	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound = errors.New("no record found")
)

type Store struct {
	Account  models.AccountStore
	Transfer models.TransferStore
}

func NewStore(db *sqlx.DB) Store {
	return Store{
		Account: AccountStore{
			db: db,
		},
		Transfer: TransferStore{
			db: db,
		},
	}
}

func NewMockStore() Store {
	return Store{
		Account: mock.MockAccountStore{},
	}
}
