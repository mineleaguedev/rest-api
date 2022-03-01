package general

import (
	"database/sql"
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
)

var (
	ErrMissingAuthValues = errors.New("missing auth values")
	ErrMissingRegValues  = errors.New("missing reg values")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidCaptcha    = errors.New("invalid captcha")
	ErrUserAlreadyExists = errors.New("username or email already exists")
	ErrUserDoesNotExist  = errors.New("user does not exist")
	ErrHashingPassword   = errors.New("error hashing password")
	ErrUnhashingPassword = errors.New("error unhashing password")

	ErrRegUser = errors.New("error registering user")

	ErrCountingUserSessions = errors.New("error counting user sessions")
	ErrAddingSessionInfo    = errors.New("error adding session information to database")
	ErrDeletingSessionInfo  = errors.New("error deleting session information from database")

	ErrExtractingClaims = errors.New("error extracting claims")

	ErrClosingMysqlRows = errors.New("error closing mysql rows")
	ErrUnknown          = errors.New("unknown error")

	DB            *sql.DB
	JWTMiddleware *jwt.GinJWTMiddleware
)

type User struct {
	ID       float64 `json:"id"`
	Username string  `json:"username"`
}

func Setup(generalDB *sql.DB, ginJWTMiddleware *jwt.GinJWTMiddleware) {
	DB = generalDB
	JWTMiddleware = ginJWTMiddleware
}

func handleErr(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, models.Error{
		Success: false,
		Message: err.Error(),
	})
}
