package main

import (
	"log"

	"github.com/gfmanica/splitz-backend/cmd/api"
	"github.com/gfmanica/splitz-backend/config"
	"github.com/gfmanica/splitz-backend/db"
)

func main() {
	db, err := db.NewPostgreSqlStorage(config.Envs.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
