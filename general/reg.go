package general

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
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

func ConfirmRegHandler(c *gin.Context) {
	token := c.Param("token")

	regInfo, err := getRegSession(token)
	if err != nil {
		handleInternalErr(c, http.StatusInternalServerError, ErrGettingRegSession, err)
		return
	}

	if deleted, err := deleteSession(token); err != nil || deleted == 0 {
		handleInternalErr(c, http.StatusInternalServerError, ErrDeletingSession, err)
		return
	}

	if _, err := DB.Exec("INSERT INTO `users` (`username`, `email`, `password_hash`) VALUES (?, ?, ?)",
		regInfo.Username, regInfo.Email, regInfo.HashedPassword); err != nil {
		handleInternalErr(c, http.StatusInternalServerError, ErrRegUser, err)
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
		return http.StatusBadRequest, ErrMissingRegValues
	}

	if !validUsername(input.Username) {
		return http.StatusBadRequest, ErrInvalidUsername
	}

	if sevenOrMore, number := validPassword(input.Password); !sevenOrMore || !number {
		return http.StatusBadRequest, ErrInvalidPassword
	}

	if response := captchaClient.VerifyToken(input.Captcha); !response.Success {
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

	regToken := generateToken(40)

	if err := saveRegSession(regToken, input.Username, input.Email, hashedPassword); err != nil {
		log.Printf(ErrSavingRegSession.Error()+": %s\n", err.Error())
		return http.StatusInternalServerError, ErrSavingRegSession
	}

	if err := sendRegEmail(input.Email, regToken); err != nil {
		log.Printf(ErrSendingEmail.Error()+": %s\n", err.Error())
		return http.StatusInternalServerError, ErrSendingEmail
	}

	return http.StatusOK, nil
}
