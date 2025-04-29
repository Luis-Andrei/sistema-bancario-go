package main

import (
	"log"
	"net/http"
	"time"

	"bank-server/db"
	"bank-server/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar arquivo .env:", err)
	}

	// Conecta ao banco de dados
	if err := db.Connect(); err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer db.DB.Close()

	// Configura o roteador
	r := mux.NewRouter()

	// Configura as rotas
	r.HandleFunc("/accounts", handlers.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{id:[0-9]+}", handlers.GetAccount).Methods("GET")
	r.HandleFunc("/accounts/{id:[0-9]+}/deposit", handlers.Deposit).Methods("POST")
	r.HandleFunc("/accounts/{id:[0-9]+}/withdraw", handlers.Withdraw).Methods("POST")
	r.HandleFunc("/transfer", handlers.Transfer).Methods("POST")

	// Aplica o middleware de logging
	handler := loggingMiddleware(r)

	// Configura o servidor
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Servidor iniciado na porta 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
