package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yellhtet-ux/todo-golang-api-proj/env"
)


func main () {

	ctx := context.Background()

	// Config
	cfg := config {
		addr: ":8080",
		db: dbConfig{
		dsn: env.GetString("GOOSE_DBSTRING","host=localhost user=postgres password=postgres dbname=todos sslmode=disable"),
		},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout,nil))
	slog.SetDefault(logger)

	// Database
	conn, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)		
	}
	defer conn.Close()

	logger.Info("connected to database","dsn",cfg.db.dsn)

	// Application
	api := application {
	config: cfg,
	db: conn,
	}

	// Run the server
	if err := api.run(api.mount()); err != nil {
		log.Printf("Server is failed to start, err %s",err)
		os.Exit(1)
	}
}