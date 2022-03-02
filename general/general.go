package general

import (
	"database/sql"
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/hcaptcha"
	"github.com/mineleaguedev/rest-api/models"
	"html/template"
)

var (
	ErrMissingAuthValues     = errors.New("missing auth values")
	ErrMissingRegValues      = errors.New("missing reg values")
	ErrInvalidUsername       = errors.New("invalid username")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidCaptcha        = errors.New("invalid captcha")
	ErrUserAlreadyExists     = errors.New("username or email already exists")
	ErrUserDoesNotExist      = errors.New("user does not exist")
	ErrHashingPassword       = errors.New("error hashing password")
	ErrUnhashingPassword     = errors.New("error unhashing password")
	ErrRegUser               = errors.New("error registering user")
	ErrCountingUserSessions  = errors.New("error counting user sessions")
	ErrAddingSessionInfo     = errors.New("error adding session information to database")
	ErrDeletingSessionInfo   = errors.New("error deleting session information from database")
	ErrInvalidToken          = errors.New("invalid token")
	ErrChangedClientIp       = errors.New("changed client ip")
	ErrNotExpiredAccessToken = errors.New("access token is not expired")
	ErrExpiredRefreshToken   = errors.New("refresh token is expired")
	ErrExtractingClaims      = errors.New("error extracting claims")
	ErrUnknown               = errors.New("unknown error")

	DB            *sql.DB
	JWTMiddleware *jwt.GinJWTMiddleware
)

type User struct {
	ID       float64 `json:"id"`
	Username string  `json:"username"`
}

func Setup(generalDB *sql.DB, ginJWTMiddleware *jwt.GinJWTMiddleware, siteKey, secretKey string) {
	DB = generalDB
	JWTMiddleware = ginJWTMiddleware

	// configure captcha
	SiteKey = siteKey
	CaptchaClient = hcaptcha.New(secretKey)
	RegForm = template.Must(template.ParseFiles("./forms/reg_form.html"))
	AuthForm = template.Must(template.ParseFiles("./forms/auth_form.html"))
}

func handleErr(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, models.Error{
		Success: false,
		Message: err.Error(),
	})
}
