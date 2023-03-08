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

var UpdateAccount = func(acc *models.Account) error {
	return nil
}

var UpdateBalance = func(acc *models.Account) error {
	return nil
}

var ListAccounts = func() ([]*models.Account, error) {
	return nil, nil
}

func (mockStore MockAccountStore) Create(acc *models.Account) error {
	return CreateAccount(acc)
}

func (mockStore MockAccountStore) Get(id int64) (*models.Account, error) {
	return GetAccountByID(id)
}
func (mockStore MockAccountStore) UpdateAccount(acc *models.Account) error {
	return UpdateAccount(acc)
}

func (mockStore MockAccountStore) UpdateBalance(acc *models.Account) error {
	return UpdateBalance(acc)
}

func (mockStore MockAccountStore) List() ([]*models.Account, error) {
	return ListAccounts()
}
