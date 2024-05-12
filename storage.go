package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccount(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStorage) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateAccount(account *Account) error {

	query := `
		INSERT INTO account
		(first_name, last_name, number, balance, created_at)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING id
	`

	var accountID int64
	err := s.db.QueryRow(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt).Scan(&accountID)
	if err != nil {
		return err
	}

	account.ID = accountID

	return nil
}

func (s *PostgresStorage) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE id = $1`
	_, err := s.db.Query(query, id)
	return err
}

func (s *PostgresStorage) GetAccount(id int) (*Account, error) {

	query := `
		SELECT * FROM account WHERE id = $1
	`
	account := new(Account)
	err := s.db.QueryRow(query, id).Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PostgresStorage) GetAllAccounts() ([]*Account, error) {
	query := `
		SELECT * FROM account
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {

		account := new(Account)
		err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}

	return accounts, nil
}
