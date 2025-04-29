package handlers

import (
	"encoding/json"
	"net/http"

	"bank-server/db"
	"bank-server/models"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc models.Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		writeError(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	if acc.Name == "" {
		writeError(w, "Nome é obrigatório", http.StatusBadRequest)
		return
	}

	if acc.Balance < 0 {
		writeError(w, "Saldo inicial não pode ser negativo", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING id`
	err := db.DB.QueryRow(query, acc.Name, acc.Balance).Scan(&acc.ID)
	if err != nil {
		writeError(w, "Erro ao criar conta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(acc)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var acc models.Account
	err := db.DB.Get(&acc, "SELECT * FROM accounts WHERE id=$1", id)
	if err != nil {
		writeError(w, "Conta não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func Deposit(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct{ Amount float64 }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		writeError(w, "Valor do depósito deve ser maior que zero", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", req.Amount, id)
	if err != nil {
		writeError(w, "Erro ao realizar depósito: "+err.Error(), http.StatusInternalServerError)
		return
	}

	GetAccount(w, r)
}

func Withdraw(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct{ Amount float64 }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		writeError(w, "Valor do saque deve ser maior que zero", http.StatusBadRequest)
		return
	}

	// Verifica saldo atual
	var current float64
	err := db.DB.Get(&current, "SELECT balance FROM accounts WHERE id = $1", id)
	if err != nil {
		writeError(w, "Conta não encontrada", http.StatusNotFound)
		return
	}

	if current < req.Amount {
		writeError(w, "Saldo insuficiente", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", req.Amount, id)
	if err != nil {
		writeError(w, "Erro ao realizar saque: "+err.Error(), http.StatusInternalServerError)
		return
	}

	GetAccount(w, r)
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	var transfer struct {
		FromID int     `json:"from_id"`
		ToID   int     `json:"to_id"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		writeError(w, "Erro ao decodificar requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	if transfer.Amount <= 0 {
		writeError(w, "Valor da transferência deve ser maior que zero", http.StatusBadRequest)
		return
	}

	// Verifica saldo da conta de origem
	var fromBalance float64
	err := db.DB.Get(&fromBalance, "SELECT balance FROM accounts WHERE id = $1", transfer.FromID)
	if err != nil {
		writeError(w, "Conta de origem não encontrada", http.StatusNotFound)
		return
	}

	if fromBalance < transfer.Amount {
		writeError(w, "Saldo insuficiente na conta de origem", http.StatusBadRequest)
		return
	}

	// Verifica se conta de destino existe
	var toExists bool
	err = db.DB.Get(&toExists, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", transfer.ToID)
	if err != nil || !toExists {
		writeError(w, "Conta de destino não encontrada", http.StatusNotFound)
		return
	}

	// Inicia transação
	tx, err := db.DB.Beginx()
	if err != nil {
		writeError(w, "Erro ao iniciar transação: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Atualiza saldo da conta de origem
	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", transfer.Amount, transfer.FromID)
	if err != nil {
		tx.Rollback()
		writeError(w, "Erro ao debitar conta de origem: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Atualiza saldo da conta de destino
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", transfer.Amount, transfer.ToID)
	if err != nil {
		tx.Rollback()
		writeError(w, "Erro ao creditar conta de destino: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Confirma transação
	if err := tx.Commit(); err != nil {
		writeError(w, "Erro ao confirmar transação: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transferência realizada com sucesso"})
}
