package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type DbInfo struct {
	User     string
	Password string
	DbName   string
	Host     string
	Port     int
}

func (dbInfo DbInfo) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbInfo.Host,
		dbInfo.User,
		dbInfo.Password,
		dbInfo.DbName,
		dbInfo.Port,
	)
}

func GetDBConnectionString() string {
	envSrcUser := "POSTGRES_USER"
	envSrcPass := "POSTGRES_PASSWORD"
	envSrcDbName := "POSTGRES_DB"
	envSrcPort := "POSTGRES_PORT"
	envSrcHost := "POSTGRES_HOST"

	port, err := strconv.Atoi(GetValueFromEnv(envSrcPort, "5432"))
	if err != nil {
		log.Fatal("Failed to convert port to int")
	}
	var dbInfo DbInfo

	dbInfo.User = GetValueFromEnv(envSrcUser, "test")
	dbInfo.Password = GetValueFromEnv(envSrcPass, "test")
	dbInfo.DbName = GetValueFromEnv(envSrcDbName, "test")
	dbInfo.Host = GetValueFromEnv(envSrcHost, "localhost")
	dbInfo.Port = port

	return dbInfo.ConnectionString()
}

func GetValueFromEnv(envSrc string, defaultValue string) string {
	envValue := os.Getenv(envSrc)
	if envValue == "" {
		return defaultValue
	}

	return envValue
}
