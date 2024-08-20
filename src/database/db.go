package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("mysql", "root:lobio2541@tcp(127.0.0.1:3306)/biblia")
	if err != nil {
		log.Fatal(err)
	}
}
