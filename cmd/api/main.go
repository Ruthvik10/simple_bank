package main

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	db, err := initDB(dsn)
	if err != nil {
		log.Fatal("error connecting to the database")
	}
	defer db.Close()
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
