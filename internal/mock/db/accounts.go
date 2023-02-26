package mock

import "github.com/Ruthvik10/simple_bank/internal/models"

type MockAccountStore struct {
}

var CreateAccount = func(acc *models.Account) error {
	return nil
}

var GetAccountByID = func(id int64) (*models.Account, error) {
	return nil, nil
}

func (mockStore MockAccountStore) Create(acc *models.Account) error {
	return CreateAccount(acc)
}

func (mockStore MockAccountStore) Get(id int64) (*models.Account, error) {
	return GetAccountByID(id)
}
