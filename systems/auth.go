package systems

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
)

type User struct {
	Username string `json:"username"`
}

func HelloHandler(c *gin.Context) {
	user, _ := c.Get("username")
	c.JSON(200, gin.H{
		"userName": user.(*User).Username,
	})
}

func Payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			"username": v.Username,
		}
	}
	return jwt.MapClaims{}
}

func Identify(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		Username: claims["username"].(string),
	}
}

func Authenticate(c *gin.Context) (interface{}, error) {
	var input models.AuthRequest

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return nil, jwt.ErrMissingLoginValues
	}

	response := CaptchaClient.VerifyToken(input.Captcha)
	if !response.Success {
		return nil, ErrInvalidCaptcha
	}

	var hashedPassword string
	err := GeneralDB.QueryRow("SELECT `password_hash` FROM `users` WHERE `username` = ?", input.Username).Scan(&hashedPassword)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrUserDoesNotExist
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, hashedPassword)
	if err != nil {
		return nil, ErrUnhashingPassword
	}

	if !match {
		c.JSON(http.StatusOK, models.Response{
			Success: false,
		})
		return nil, ErrUnknown
	}

	return &User{
		Username: input.Username,
	}, nil
}

func Authorize(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok && v.Username != "" {
		return true
	}
	return false
}
