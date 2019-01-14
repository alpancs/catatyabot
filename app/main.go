package main

import (
	"database/sql"
	"net/http"
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

func main() {
	http.HandleFunc("/"+os.Getenv("BOT_TOKEN"), handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
