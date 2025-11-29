package main

import (
	"context"
	"ecomApis/internals/env"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	appconfig := appconfig{
		Address: env.GetString("APP_ADDRESS", ":8080"),
		DB: dbConfig{
			DatabaseURL: env.GetString("GOOSE_DBSTRING", "host=localhost port=5433 user=postgres password=postgres dbname=ecommerce_db sslmode=disable"),
		},
	}

	// use slog for structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// setup the db
	conn, err := pgx.Connect(ctx, appconfig.DB.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to the database successfully")

	app := &application{
		config: appconfig,
		db:     conn,
	}

	// start the server
	err = app.run(app.mount())
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}

}
