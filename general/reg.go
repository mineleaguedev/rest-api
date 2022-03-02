package general

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
	"regexp"
	"unicode"
)

func RegHandler(c *gin.Context) {
	httpCode, err := registerUser(c)
	if err != nil {
		handleErr(c, httpCode, err)
		return
	}

	c.JSON(httpCode, models.Response{
		Success: true,
	})
}

func validUsername(username string) bool {
	matched, err := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", username)
	if err != nil {
		return false
	}
	return matched
}

func validPassword(password string) (sevenOrMore, number bool) {
	letters := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			letters++
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	sevenOrMore = letters >= 7
	return
}

func registerUser(c *gin.Context) (int, error) {
	var input models.RegisterRequest

	if err := c.ShouldBind(&input); err != nil {
		return http.StatusBadRequest, ErrMissingRegValues
	}

	if !validUsername(input.Username) {
		return http.StatusBadRequest, ErrInvalidUsername
	}

	if sevenOrMore, number := validPassword(input.Password); !sevenOrMore || !number {
		return http.StatusBadRequest, ErrInvalidPassword
	}

	if response := CaptchaClient.VerifyToken(input.Captcha); !response.Success {
		return http.StatusBadRequest, ErrInvalidCaptcha
	}

	var exists bool
	err := DB.QueryRow("SELECT 1 FROM `users` WHERE `username` = ? OR `email` = ?", input.Username, input.Email).Scan(&exists)
	if (err != nil && err != sql.ErrNoRows) || exists {
		return http.StatusBadRequest, ErrUserAlreadyExists
	}

	hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		log.Printf(ErrHashingPassword.Error()+": %s\n", err.Error())
		return http.StatusInternalServerError, ErrHashingPassword
	}

	if _, err := DB.Exec("INSERT INTO `users` (`username`, `email`, `password_hash`) VALUES (?, ?, ?)",
		input.Username, input.Email, hashedPassword); err != nil {
		log.Printf(ErrRegUser.Error()+": %s\n", err.Error())
		return http.StatusInternalServerError, ErrRegUser
	}

	return http.StatusOK, nil
}
