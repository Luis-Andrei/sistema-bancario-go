package config

import (
	"database/sql"
	"fmt"
	"log"

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

func GetDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
