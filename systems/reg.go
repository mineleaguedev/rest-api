package systems

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

func validUsername(username string) bool {
	matched, err := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", username)
	if err != nil {
		return false
	}
	return matched
}

func validPassword(password string) (sevenOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	sevenOrMore = letters >= 7
	return
}

func RegisterUser(c *gin.Context) {
	var input models.RegisterRequest

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return
	}

	if !validUsername(input.Username) {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Invalid username",
		})
		return
	}

	var sevenOrMore, number, upper, special = validPassword(input.Password)
	if !sevenOrMore || !number || !upper || !special {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Invalid password",
		})
		return
	}

	response := CaptchaClient.VerifyToken(input.Captcha)
	if !response.Success {
		c.JSON(http.StatusUnauthorized, models.Error{
			Success: false,
			Message: "Invalid captcha",
		})
		return
	}

	var exists bool
	err := GeneralDB.QueryRow("SELECT 1 FROM `users` WHERE `username` = ? OR `email` = ?", input.Username, input.Email).Scan(&exists)
	if (err != nil && err != sql.ErrNoRows) || exists {
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Username or email already exists",
		})
		return
	}

	hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error hashing password",
		})
		log.Printf("Error hashing password: %s\n", err.Error())
		return
	}

	if _, err := GeneralDB.Exec("INSERT INTO `users` (`username`, `email`, `password_hash`) VALUES (?, ?, ?)",
		input.Username, input.Email, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error registering user",
		})
		log.Printf("Error registering user: %s\n", err.Error())
	} else {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
		})
	}
}
