package admin

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/services"
)

type Handler struct {
	services    *services.Service
	generalDB   *sql.DB
	minigamesDB *sql.DB
}

func NewHandler(services *services.Service, generalDB, minigamesDB *sql.DB) *Handler {
	return &Handler{
		services:    services,
		generalDB:   generalDB,
		minigamesDB: minigamesDB,
	}
}
