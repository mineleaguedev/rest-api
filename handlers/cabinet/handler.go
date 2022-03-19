package cabinet

import (
	"database/sql"
	"github.com/mineleaguedev/rest-api/services"
)

type Handler struct {
	services *services.Service
	db       *sql.DB
}

var imageTypes = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

func NewHandler(services *services.Service, db *sql.DB) *Handler {
	return &Handler{
		services: services,
		db:       db,
	}
}
