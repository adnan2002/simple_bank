package main

import (
	"context"

	"example.com/api"
	"example.com/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dbSource = "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable"

func main() {
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		panic("Couldn't connect to db")
	}
	dbStore := db.NewStore(dbPool)

	srv := api.NewServer(dbStore)

	if err := srv.Start(":8032"); err != nil {
		panic("Couldn't run server")
	}

}
