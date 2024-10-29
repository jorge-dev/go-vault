package main

import "math/rand"

type Account struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	LastName string  `json:"last_name"`
	Number   int64   `json:"number"`
	Balance  float64 `json:"balance"`
}

func NewAccount(id int, name, lastName, number string, balance float64) *Account {
	return &Account{
		Id:       rand.Intn(1000),
		Name:     name,
		LastName: lastName,
		Number:   rand.Int63n(1000000000000000),
		Balance:  balance,
	}
}
