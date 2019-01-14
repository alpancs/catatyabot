package app

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var db = newDatabaseConnection()

func newDatabaseConnection() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return db
}
