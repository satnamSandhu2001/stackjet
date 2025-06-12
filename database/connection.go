package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/satnamSandhu2001/stackjet/pkg"
	_ "modernc.org/sqlite"
)

func Connect() *sqlx.DB {
	db, err := sqlx.Open("sqlite", pkg.Config().DB_URL)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	return db
}
