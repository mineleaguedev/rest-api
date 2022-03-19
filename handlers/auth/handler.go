package auth

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
	"github.com/mineleaguedev/rest-api/services"
)

type Handler struct {
	services   *services.Service
	middleware models.JWTMiddleware
	db         *sql.DB
}

func NewHandler(services *services.Service, middleware models.JWTMiddleware, db *sql.DB) *Handler {
	return &Handler{
		services:   services,
		middleware: middleware,
		db:         db,
	}
}
