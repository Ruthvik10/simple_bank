package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ruthvik10/simple_bank/internal/logger"
	"github.com/Ruthvik10/simple_bank/internal/store"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	cfg    config
	logger *logger.Logger
	store  store.Store
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cfg := config{
		port: os.Getenv("PORT"),
		env:  os.Getenv("ENV"),
		db: struct{ dsn string }{
			dsn: os.Getenv("DSN"),
		},
	}
	l := logger.New(os.Stdout, logger.LevelInfo)
	app := application{
		cfg:    cfg,
		logger: l,
	}
	db, err := initDB(app.cfg.db.dsn)
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
	defer db.Close()
	app.store = store.NewStore(db, l)
	app.logger.PrintInfo("database connection successful", nil)
	app.logger.PrintInfo("starting server on port "+app.cfg.port, nil)
	err = initServer(app.cfg.port, app.routes(), l)
	if err != nil {
		app.logger.PrintError(err, nil)
	}
}

func initServer(port string, handler http.Handler, logger *logger.Logger) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     log.New(logger, "", 0),
	}
	return srv.ListenAndServe()
}

func initDB(dsn string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
