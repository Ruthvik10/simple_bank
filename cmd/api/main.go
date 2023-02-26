package main

import (
	"context"
	"log"
	"os"
	"time"

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
	logger *log.Logger
	store  store.Store
}

func init() {
	err := godotenv.Load("../../.env")
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
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := application{
		cfg:    cfg,
		logger: l,
	}
	db, err := initDB(app.cfg.db.dsn)
	if err != nil {
		log.Fatal("error connecting to the database", err)
	}
	defer db.Close()
	app.store = store.NewStore(db)
	log.Println("database connection successful")
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
