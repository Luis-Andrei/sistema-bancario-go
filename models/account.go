package models

import (
	"database/sql"
	"errors"
)

type Account struct {
	ID      int     `db:"id" json:"id"`
	Name    string  `db:"name" json:"name"`
	Balance float64 `db:"balance" json:"balance"`
}

func CreateAccount(db *sql.DB, account *Account) error {
	query := `INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING id`
	return db.QueryRow(query, account.Name, account.Balance).Scan(&account.ID)
}

func GetAccount(db *sql.DB, id int) (*Account, error) {
	account := &Account{}
	query := `SELECT id, name, balance FROM accounts WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&account.ID, &account.Name, &account.Balance)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func UpdateBalance(db *sql.DB, id int, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	result, err := db.Exec(query, amount, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("conta não encontrada")
	}

	return nil
}

func Transfer(db *sql.DB, fromID, toID int, amount float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Verificar saldo da conta de origem
	var balance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = $1", fromID).Scan(&balance)
	if err != nil {
		tx.Rollback()
		return err
	}

	if balance < amount {
		tx.Rollback()
		return errors.New("saldo insuficiente")
	}

	// Atualizar saldo da conta de origem
	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Atualizar saldo da conta de destino
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
