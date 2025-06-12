package database

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/satnamSandhu2001/stackjet/pkg"
)

func RunInitSQL() error {
	conn := Connect()
	defer conn.Close()

	_, err := conn.Exec(string(InitSQL))
	if err != nil {
		return fmt.Errorf("init.sql execution failed: %w", err)
	}

	// Insert default admin
	email := "admin@stackjet.com"
	password := "admin123"

	hashed, err := pkg.GenerateHash(password)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`
		INSERT INTO users (email, password, role)
		VALUES (?, ?, 'superadmin')
	`, email, hashed)

	if err != nil {
		return fmt.Errorf("admin insert failed: %w", err)
	}

	fmt.Println("\033[1;34m🔐 Admin user created with the following credentials:\033[0m")
	fmt.Printf("\033[34m      email:    %s\033[0m\n", email)
	fmt.Printf("\033[34m      password: %s\033[0m\n\n", password)

	return nil
}
