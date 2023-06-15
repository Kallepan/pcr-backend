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

func GetConnectionString() string {
	dbInfo := getDbInfo()

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbInfo.Host,
		dbInfo.User,
		dbInfo.Password,
		dbInfo.DbName,
		dbInfo.Port,
	)

	return connectionString
}

func getDbInfo() DbInfo {
	envSrcUser := "POSTGRES_USER"
	envSrcPass := "POSTGRES_PASSWORD"
	envSrcDbName := "POSTGRES_DB"
	envSrcPort := "POSTGRES_PORT"
	envSrcHost := "POSTGRES_HOST"

	port, err := strconv.Atoi(GetValueFromEnv(envSrcPort, "5432"))
	if err != nil {
		log.Fatal("Failed to convert port to int")
	}

	dbInfo := DbInfo{
		User:     GetValueFromEnv(envSrcUser, "test"),
		Password: GetValueFromEnv(envSrcPass, "test"),
		DbName:   GetValueFromEnv(envSrcDbName, "test"),
		Host:     GetValueFromEnv(envSrcHost, "localhost"),
		Port:     port,
	}

	return dbInfo
}

func GetValueFromEnv(envSrc string, defaultValue string) string {
	envValue := os.Getenv(envSrc)
	if envValue == "" {
		return defaultValue
	}

	return envValue
}
