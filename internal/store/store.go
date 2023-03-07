package store

import (
	"errors"

	"github.com/Ruthvik10/simple_bank/internal/logger"
	mock "github.com/Ruthvik10/simple_bank/internal/mock/db"
	"github.com/Ruthvik10/simple_bank/internal/models"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound = errors.New("no record found")
)

type Store struct {
	Account models.AccountStore
}

func NewStore(db *sqlx.DB, l *logger.Logger) Store {
	return Store{
		Account: AccountStore{
			db:     db,
			logger: l,
		},
	}
}

func NewMockStore() Store {
	return Store{
		Account: mock.MockAccountStore{},
	}
}
