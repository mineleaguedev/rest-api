package minigames

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/services"
)

type Handler struct {
	services *services.Service
	db       *sql.DB
}

func NewHandler(services *services.Service, db *sql.DB) *Handler {
	return &Handler{
		services: services,
		db:       db,
	}
}
