package controllers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var GeneralDB *sql.DB
var MiniGamesDB *sql.DB

func Controller(generalDB *sql.DB, miniGamesDB *sql.DB) {
	GeneralDB = generalDB
	MiniGamesDB = miniGamesDB
}
