package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/theborzet/time-tracker/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init(cfg *config.Config) *sql.DB {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
