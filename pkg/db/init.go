package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(connection string) *sql.DB {
	db, err := sql.Open("mysql", connection)
    if err != nil {
        panic(err.Error())
    }
    if err := db.Ping(); err != nil {
        panic(err.Error())
    }
    return db
}
