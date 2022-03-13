package general

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
	"regexp"
	"unicode"
)

func RegHandler(c *gin.Context) {
	httpCode, err := registerUser(c)
	if err != nil {
		Service.HandleErr(c, httpCode, err)
		return
	}

	c.JSON(httpCode, models.Response{
		Success: true,
	})
}

func ConfirmRegHandler(c *gin.Context) {
	token := c.Param("token")

	regInfo, err := Service.GetRegSession(token)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrGettingRegSession, err)
		return
	}

	if deleted, err := Service.DeleteSession(token); err != nil || deleted == 0 {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrDeletingSession, err)
		return
	}

	if _, err := DB.Exec("INSERT INTO `users` (`username`, `email`, `password_hash`) VALUES (?, ?, ?)",
		regInfo.Username, regInfo.Email, regInfo.HashedPassword); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrRegisteringUser, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
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

func generateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func registerUser(c *gin.Context) (int, error) {
	var input models.RegisterRequest

	if err := c.ShouldBind(&input); err != nil {
		return http.StatusBadRequest, errors.ErrMissingRegValues
	}

	if !validUsername(input.Username) {
		return http.StatusBadRequest, errors.ErrInvalidUsername
	}

	if sevenOrMore, number := validPassword(input.Password); !sevenOrMore || !number {
		return http.StatusBadRequest, errors.ErrInvalidPassword
	}

	if response := Service.VerifyCaptcha(input.Captcha); !response.Success {
		return http.StatusBadRequest, errors.ErrInvalidCaptcha
	}

	var exists bool
	err := DB.QueryRow("SELECT 1 FROM `users` WHERE `username` = ? OR `email` = ?", input.Username, input.Email).Scan(&exists)
	if (err != nil && err != sql.ErrNoRows) || exists {
		return http.StatusBadRequest, errors.ErrUserAlreadyExists
	}

	hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		log.Printf(errors.ErrHashingPassword.Error()+": %s\n", err.Error())
		return http.StatusInternalServerError, errors.ErrHashingPassword
	}

	regToken := generateToken(40)

	if err := Service.SaveRegSession(regToken, input.Username, input.Email, hashedPassword, Middleware.RegTokenTime); err != nil {
		log.Printf(errors.ErrSavingRegSession.Error()+": %s\n", err.Error())
		return http.StatusInternalServerError, errors.ErrSavingRegSession
	}

	if err := Service.SendRegEmail(input.Email, regToken); err != nil {
		log.Printf(errors.ErrSendingEmail.Error()+": %s\n", err.Error())
		return http.StatusInternalServerError, errors.ErrSendingEmail
	}

	return http.StatusOK, nil
}
