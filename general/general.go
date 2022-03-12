package general

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
	"github.com/mineleaguedev/rest-api/services"
)

var (
	DB         *sql.DB
	Service    *services.Service
	Middleware models.JWTMiddleware
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

func Setup(generalDB *sql.DB, service *services.Service, jwtMiddleware models.JWTMiddleware) {
	DB = generalDB
	Service = service
	Middleware = jwtMiddleware
}
