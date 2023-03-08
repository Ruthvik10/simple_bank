package models

import (
	"time"
)

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   int64     `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AccountStore interface {
	Get(int64) (*Account, error)
	List() ([]*Account, error)
	Create(*Account) error
	UpdateAccount(*Account) error
	UpdateBalance(*Account) error
}
