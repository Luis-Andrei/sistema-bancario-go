package models

import (
	"database/sql"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	Number    string    `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateAccount(db *sql.DB, account *Account) error {
	query := `
		INSERT INTO accounts (number, balance, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	return db.QueryRow(query, account.Number, account.Balance, time.Now()).Scan(&account.ID)
}

func GetAccount(db *sql.DB, id int) (*Account, error) {
	account := &Account{}
	query := `SELECT id, number, balance, created_at FROM accounts WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&account.ID, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func UpdateBalance(db *sql.DB, id int, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := db.Exec(query, amount, id)
	return err
}
