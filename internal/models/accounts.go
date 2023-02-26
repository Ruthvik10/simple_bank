package models

import (
	"time"
)

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   string    `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Version   int64     `json:"version" db:"version"`
}

type AccountStore interface {
	Get(int64) (*Account, error)
	// List() ([]*Account, error)
	Create(*Account) error
	// UpdateBalance(*Account) error
}
