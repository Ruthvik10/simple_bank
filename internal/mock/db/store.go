package mock

type MockStore struct {
	Account MockAccountStore
}

func NewMockStore() MockStore {
	return MockStore{
		Account: MockAccountStore{},
	}
}
