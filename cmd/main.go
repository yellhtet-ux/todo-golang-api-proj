package main

import (
	"log"
	"log/slog"
	"os"
)


func main () {
	cfg := config {
		addr: ":8000",
		db: dbConfig{},
	}

	app := application {
		config: cfg,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout,nil))
	slog.SetDefault(logger)

	if err := app.run(app.mount()); err != nil {
		log.Printf("Server is failed to start, err %s",err)
		os.Exit(1)
	}
}