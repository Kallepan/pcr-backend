package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func Migrate() {
	if Instance == nil {
		panic(errors.New("database instance is not initialized"))
	}

	driver, err := postgres.WithInstance(Instance, &postgres.Config{})

	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	log.Println("Migrations completed!")
}
