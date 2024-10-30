package main

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}

type Account struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Number    uuid.UUID `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(name, lastName string, balance float64) *Account {
	return &Account{
		Name:      name,
		LastName:  lastName,
		Number:    uuid.New(),
		Balance:   balance,
		CreatedAt: time.Now().UTC(),
	}
}
