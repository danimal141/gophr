package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var globalMySQLDB *sql.DB

func init() {
	db, err := NewMySQLDB(os.Getenv("MYSQL_DATA_SOURCE_NAME"))
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db
}

func NewMySQLDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
