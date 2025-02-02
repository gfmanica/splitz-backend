package main

import (
	"context"
	"log"

	"github.com/gfmanica/splitz-backend/cmd/api"
	"github.com/gfmanica/splitz-backend/config"
	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), config.Envs.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", conn)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
