package controllers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Controller(db *sql.DB) {
	DB = db
}
