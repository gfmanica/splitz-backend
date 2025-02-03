package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgreSqlStorage(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseUrl)

	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
