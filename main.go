package main

import (
	"log"
	"net/http"
	"projeto/config"
	"projeto/handlers"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	// Inicializa o banco de dados
	db := config.InitDB()
	defer db.Close()

	// Cria a tabela de contas se não existir
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			number VARCHAR(20) UNIQUE NOT NULL,
			balance DECIMAL(10,2) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Inicializa o handler
	accountHandler := &handlers.AccountHandler{DB: db}

	// Configura o router
	router := mux.NewRouter()

	// Configura as rotas
	router.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}", accountHandler.GetAccount).Methods("GET")
	router.HandleFunc("/accounts/{id}/deposit", accountHandler.Deposit).Methods("POST")
	router.HandleFunc("/accounts/{id}/withdraw", accountHandler.Withdraw).Methods("POST")

	// Configura a proteção CSRF
	CSRF := csrf.Protect(
		[]byte("32-byte-long-auth-key"), // Chave secreta para o CSRF
		csrf.Secure(false),              // Desativa em desenvolvimento
	)

	// Inicia o servidor
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", CSRF(router)))
}
