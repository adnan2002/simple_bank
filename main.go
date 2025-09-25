package main

import (
	"context"
	"fmt"
	"log"

	"example.com/api"
	"example.com/db/sqlc"
	"example.com/db/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	config, err := util.LoadConfig("app.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Build DB connection string
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbSslMode,
	)

	// Initialize DB connection pool
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Create store and server
	store := db.NewStore(dbPool)
	server := api.NewServer(store)

	// Start server
	addr := fmt.Sprintf(":%s", config.AppPort)
	log.Printf("starting server on %s...", addr)

	if err := server.Start(addr); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
