package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // postgres driver
)

var Instance *sql.DB
var dbError error

func Connect(connectionString string) {
	Instance, dbError = sql.Open("postgres", connectionString)

	if dbError != nil {
		log.Fatal("Can not connect to DB: ", dbError)
		panic("Failed to connect to database!")
	}

	dbError = Instance.Ping()

	if dbError != nil {
		log.Fatal("Can not connect to DB: ", dbError)
		panic("Failed to connect to database!")
	}

	log.Println("Connected to database!")
}
