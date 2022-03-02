package general

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/hcaptcha"
	"github.com/mineleaguedev/rest-api/models"
	"html/template"
	"time"
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
	ErrCountingUserSessions    = errors.New("error counting user sessions")
	ErrAddingSessionInfo       = errors.New("error adding session information to database")
	ErrDeletingSessionInfo     = errors.New("error deleting session information from database")
	ErrInvalidToken            = errors.New("invalid token")
	ErrChangedClientIp         = errors.New("changed client ip")
	ErrNotExpiredAccessToken   = errors.New("access token is not expired")
	ErrExpiredRefreshToken     = errors.New("refresh token is expired")
	ErrExtractingClaims        = errors.New("error extracting claims")
	ErrEmptyAuthHeader         = errors.New("auth header is empty")
	ErrInvalidAuthHeader       = errors.New("auth header is invalid")
	ErrEmptyQueryToken         = errors.New("query token is empty")
	ErrEmptyCookieToken        = errors.New("cookie token is empty")
	ErrEmptyParamToken         = errors.New("parameter token is empty")
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")
	ErrMissingExpField         = errors.New("missing exp field")
	ErrWrongFormatOfExp        = errors.New("exp must be float64 format")
	ErrExpiredToken            = errors.New("token is expired")
	ErrForbidden               = errors.New("you don't have permission to access this resource")
	ErrFailedTokenCreation     = errors.New("failed to create JWT Token")

	DB         *sql.DB
	Middleware models.JWTMiddleware
)

type User struct {
	ID       float64 `json:"id"`
	Username string  `json:"username"`
}

func Setup(generalDB *sql.DB, middleware models.JWTMiddleware, siteKey, secretKey string) {
	DB = generalDB
	Middleware = middleware

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

func clientIp(c *gin.Context) string {
	clientIp := c.ClientIP()
	if clientIp == "::1" {
		clientIp = "127.0.0.1"
	}
	return clientIp
}

func createToken(mapClaims jwt.MapClaims) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(Middleware.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	for key, value := range mapClaims {
		claims[key] = value
	}

	expire := time.Now().Add(Middleware.AccessTokenTime)
	claims["exp"] = expire.Unix()
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(Middleware.Key)
	return tokenString, expire, err
}

func sendCookie(c *gin.Context, token string) {
	if Middleware.SendCookie {
		expire := time.Now().Add(Middleware.CookieMaxAge)
		maxAge := int(expire.Unix() - time.Now().Unix())

		if Middleware.CookieSameSite != 0 {
			c.SetSameSite(Middleware.CookieSameSite)
		}

		c.SetCookie(
			Middleware.CookieName,
			token,
			maxAge,
			"/",
			Middleware.CookieDomain,
			Middleware.SecureCookie,
			Middleware.CookieHTTPOnly,
		)
	}
}

func extractClaimsFromContext(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)
}
