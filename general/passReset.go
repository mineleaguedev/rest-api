package general

import (
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"math/rand"
	"net/http"
	"strings"
)

var (
	lowerCharSet = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberSet    = "0123456789"
	allCharSet   = lowerCharSet + upperCharSet + numberSet
)

func PassResetHandler(c *gin.Context) {
	var input models.PassResetRequest

	if err := c.ShouldBind(&input); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPassResetValues)
		return
	}

	if response := Service.VerifyCaptcha(input.Captcha); !response.Success {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	var userId int64
	err := DB.QueryRow("SELECT `id` FROM `users` WHERE `username` = ? OR `email` = ?", input.Username, input.Email).Scan(&userId)
	if err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		return
	}

	token := generateToken(40)

	if err := Service.SavePassResetSession(token, userId, Middleware.PassResetTokenTime); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSavingPassResetSession, err)
		return
	}

	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	if err := Service.SendPassResetEmail(input.Email, token, input.Username, ip); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSendingEmail, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func generatePassword(passwordLength, minNum, minUpperCase int) string {
	var password strings.Builder

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func ConfirmPassResetHandler(c *gin.Context) {
	token := c.Param("token")

	userId, err := Service.GetPassResetSession(token)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrGettingPassResetSession, err)
		return
	}

	if deleted, err := Service.DeleteSession(token); err != nil || deleted == 0 {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrDeletingSession, err)
		return
	}

	var email string
	var username string
	if err := DB.QueryRow("SELECT `email`, `username` FROM `users` WHERE `id` = ?", userId).Scan(&email, &username); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrUserDoesNotExist, err)
		return
	}

	newPassword := generatePassword(12, 1, 1)
	hashedPassword, err := argon2id.CreateHash(newPassword, argon2id.DefaultParams)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrHashingPassword, err)
		return
	}

	if _, err := DB.Exec("UPDATE `users` SET `password_hash` = ? WHERE `id` = ?", hashedPassword, userId); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrUpdatingUserPassword, err)
		return
	}

	if err := Service.SendNewPassEmail(email, username, newPassword); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSendingEmail, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
