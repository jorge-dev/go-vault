package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccountById(id int) (*Account, error)
	GetAccounts() ([]*Account, error)
	CreateAccount(account *Account) error
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error
}

type PostgresStore struct {
	// db connection
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connectionString := "user=user dbname=bank password=password sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{DB: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			number UUID,
			balance DOUBLE PRECISION,
			created_at TIMESTAMP DEFAULT NOW()
			)`
	_, err := s.DB.Exec(query)
	return err

}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1`
	rows, err := s.DB.Query(query, id)
	if err != nil {
		log.Println("Error getting account by id: from db", err)
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found", id)

}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM accounts`
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println("Error getting accounts:", err)
		return nil, err
	}

	accounts := make([]*Account, 0)
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `
		INSERT INTO accounts (first_name, last_name, number, balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
		`
	res, err := s.DB.Query(
		query,
		account.Name,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		log.Println("Error inserting account:", err)
		return err
	}

	log.Printf("Account created: %v", account)
	defer res.Close()
	return nil

}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`
	res, err := s.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting account:", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Account %d not found", id)
	}

	log.Printf("Account %d deleted", id)
	return nil

}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	if err := rows.Scan(&account.Id, &account.Name, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt); err != nil {
		log.Println("Error scanning account:", err)
		return nil, err
	}
	return account, nil
}
