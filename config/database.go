package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Configurações do banco de dados
// ATENÇÃO: Altere estas configurações de acordo com seu ambiente PostgreSQL
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres" // Usuário padrão do PostgreSQL
	password = "postgres" // Altere para a senha que você definiu durante a instalação
	dbname   = "banco"    // Nome do banco de dados
)

func InitDB() *sql.DB {
	// Verifica se as variáveis de ambiente estão definidas
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = user
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = password
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, dbUser, dbPassword, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar conexão com o banco de dados:", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")
	return db
}
