package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	dsn := "host=localhost port=5432 user=postgres password=12345 dbname=marketplace_of_services sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB is not connected", err)
	}

	log.Println("DB is connected")
}
