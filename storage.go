package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgreStore struct {
	db *sql.DB
}

func NewPostrgreStore() (*PostgreStore, error) {
	connStr := os.Getenv("DATABASE")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreStore{db: db}, nil
}

func (s *PostgreStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgreStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		number SERIAL NOT NULL,
		encrypted_password VARCHAR(100) NOT NULL,
		balance SERIAL NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgreStore) CreateAccount(acc *Account) error {

	query := `INSERT INTO accounts (first_name, last_name, number, encrypted_password, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.EncryptedPassword, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgreStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgreStore) DeleteAccount(id int) error {

	_, err := s.db.Query(`DELETE FROM accounts WHERE id = $1`, id)

	return err
}

func (s *PostgreStore) GetAccountByNumber(number int) (*Account, error) {

	rows, err := s.db.Query(`SELECT * FROM accounts WHERE number = $1`, number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account %d not found", number)
}

func (s *PostgreStore) GetAccountByID(id int) (*Account, error) {

	rows, err := s.db.Query(`SELECT * FROM accounts WHERE id = $1`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgreStore) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query("SELECT * FROM accounts")

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {

		account, err := scanIntoAccounts(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, err

}

func scanIntoAccounts(rows *sql.Rows) (*Account, error) {

	account := new(Account)

	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.EncryptedPassword, &account.Balance, &account.CreatedAt)

	return account, err
}
