package main

import (
	"gitlab.com/kaka/pcr-backend/database"
	"gitlab.com/kaka/pcr-backend/utils"
)

func main() {
	connectionString := utils.GetConnectionString()
	database.Connect(connectionString)
	database.Migrate()
}
