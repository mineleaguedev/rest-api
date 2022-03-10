package general

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
	"log"
)

var (
	ErrMissingAuthValues       = errors.New("missing auth values")
	ErrMissingRegValues        = errors.New("missing reg values")
	ErrInvalidUsername         = errors.New("invalid username")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrInvalidCaptcha          = errors.New("invalid captcha")
	ErrUserAlreadyExists       = errors.New("username or email already exists")
	ErrUserDoesNotExist        = errors.New("user does not exist")
	ErrHashingPassword         = errors.New("error hashing password")
	ErrUnhashingPassword       = errors.New("error unhashing password")
	ErrWrongUsernameOrPassword = errors.New("wrong username or password")
	ErrRegUser                 = errors.New("error registering user")

	ErrFailedTokenCreation     = errors.New("failed to create jwt token")
	ErrSavingAuthSession       = errors.New("error saving auth session")
	ErrGettingAuthSession      = errors.New("error getting auth session")
	ErrSavingPassResetSession  = errors.New("error saving password reset session")
	ErrGettingPassResetSession = errors.New("error getting password reset session")
	ErrSavingRegSession        = errors.New("error saving reg session")
	ErrGettingRegSession       = errors.New("error getting reg session")
	ErrDeletingSession         = errors.New("error deleting session")

	ErrSendingEmail = errors.New("error sending email")

	ErrInvalidAccessToken          = errors.New("invalid access token")
	ErrExpiredAccessToken          = errors.New("access token is expired")
	ErrAccessTokenUuidNotExists    = errors.New("failed to get access token uuid")
	ErrAccessTokenUserIdNotExists  = errors.New("failed to get access token user id")
	ErrMissingRefreshToken         = errors.New("missing refresh token")
	ErrInvalidRefreshToken         = errors.New("invalid refresh token")
	ErrExpiredRefreshToken         = errors.New("refresh token is expired")
	ErrRefreshTokenUuidNotExists   = errors.New("failed to get refresh token uuid")
	ErrRefreshTokenUserIdNotExists = errors.New("failed to get refresh token user id")

	DB         *sql.DB
	Middleware models.JWTMiddleware
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     int64
}

func Setup(generalDB *sql.DB, middleware models.JWTMiddleware) {
	DB = generalDB
	Middleware = middleware
}

func handleErr(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, models.Error{
		Success: false,
		Message: err.Error(),
	})
}

func handleInternalErr(c *gin.Context, httpCode int, err, internalErr error) {
	if internalErr != nil {
		log.Printf(err.Error()+": %s\n", internalErr.Error())
	}
	handleErr(c, httpCode, err)
}
