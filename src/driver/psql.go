package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	DB   *sql.DB
	once sync.Once
)

type dbConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     string
}

func (db *dbConfig) Init(ctx context.Context) {
	once.Do(
		func() {
			connectionString := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				db.Host,
				db.User,
				db.Password,
				db.DBName,
				db.Port,
			)
			instance, err := sql.Open("postgres", connectionString)
			if err != nil {
				slog.Error(err.Error())
				panic(err)
			}

			// Close the DB connection when the context is done
			go func() {
				<-ctx.Done()
				if err := DB.Close(); err != nil {
					slog.Error(err.Error())
					panic(err)
				}
			}()

			// Ping the DB to check if the connection is alive
			if err := instance.Ping(); err != nil {
				slog.Error(err.Error())
				panic(err)
			}

			// Set the global DB instance to the local instance
			DB = instance
		},
	)
}

func initMigrations() {
	// Run the migrations
	if DB == nil {
		slog.Error("DB connection not initialized")
		panic("DB connection not initialized")
	}

	// Run the migrations
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {
		migrationPath = "file://app/migrations"
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error(err.Error())
		panic(err)
	} else if err == migrate.ErrNoChange {
		slog.Info("No migration to run")
	}

	slog.Info("Migrations completed!")
}

func Init(ctx context.Context) {
	/* Set up the connection String using the dbConfig struct */
	db := dbConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}

	// Initialize the DB connection
	db.Init(ctx)

	// Run the migrations
	initMigrations()
}
