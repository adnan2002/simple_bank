package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"example.com/api"
	"example.com/db/sqlc"
	"example.com/db/util"
	"example.com/gapi"
	"example.com/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	config, err := util.LoadConfig("app.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Build DB connection string
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)

	// Initialize DB connection pool
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Create store and server
	store := db.NewStore(dbPool)
	runGRPCServer(store, config)
}

func runGRPCServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)

	if err != nil {
		log.Fatalf("failed to start server")
	}
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	addr := fmt.Sprintf(":%s", config.GrpcPort)
	log.Printf("starting grpc server on %s...", addr)
	listner, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("server stopped with error: %v", err)
		return
	}

	err = grpcServer.Serve(listner)

	if err != nil {
		log.Fatalf("server stopped with error: %v", err)
		return
	}

}

func runGinServer(store db.Store, config util.Config) {
	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatalf("failed to start server")
	}

	// Start server
	addr := fmt.Sprintf(":%s", config.AppPort)
	log.Printf("starting server on %s...", addr)

	if err := server.Start(addr); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}

}
