package auth

import (
	"log/slog"
	"os"

	"gitlab.com/kallepan/pcr-backend/app/domain/dao"
	"gitlab.com/kallepan/pcr-backend/driver"
)

func CreateAdminUser() {
	var user dao.User

	user.Email = os.Getenv("ADMIN_EMAIL")
	user.FirstName = "admin"
	user.LastName = "admin"
	user.Username = os.Getenv("ADMIN_USERNAME")
	user.Password = os.Getenv("ADMIN_PASSWORD")
	user.HashPassword()

	query := `
		INSERT INTO users (email, firstname, lastname, username, password, is_admin)
		VALUES ($1, $2, $3, $4, $5, true)
		ON CONFLICT (username) DO NOTHING
		`
	_, err := driver.DB.Exec(
		query,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Password)
	if err != nil {
		panic(err)
	}

	slog.Info("Admin user created")
}
