package auth

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"math/rand"
	"net/http"
	"strings"
)

const (
	lowerCharSet = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberSet    = "0123456789"
	allCharSet   = lowerCharSet + upperCharSet + numberSet
)

func (h *Handler) PassResetHandler(c *gin.Context) {
	var input models.PassResetRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPassResetValues)
		return
	}

	if response := h.services.VerifyCaptcha(input.Captcha); !response.Success {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	var userId int64
	if err := h.db.QueryRow("SELECT `id` FROM `users` WHERE `username` = ? OR `email` = ?", input.Username, input.Email).Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleInternalErr(c, errors.ErrDBGettingUser, err)
		}
		return
	}

	token := generateToken(40)

	if err := h.services.SavePassResetSession(token, userId, h.middleware.PassResetTokenTime); err != nil {
		h.services.HandleInternalErr(c, errors.ErrSavingPassResetSession, err)
		return
	}

	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	if err := h.services.SendPassResetEmail(input.Email, token, input.Username, ip); err != nil {
		h.services.HandleInternalErr(c, errors.ErrSendingEmail, err)
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

func (h *Handler) PassResetConfirmHandler(c *gin.Context) {
	token := c.Param("token")

	userId, err := h.services.GetPassResetSession(token)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrGettingPassResetSession, err)
		return
	}

	if deleted, err := h.services.DeleteSession(token); err != nil || deleted == 0 {
		h.services.HandleInternalErr(c, errors.ErrDeletingSession, err)
		return
	}

	var email, username string
	if err := h.db.QueryRow("SELECT `email`, `username` FROM `users` WHERE `id` = ?", userId).Scan(&email, &username); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleInternalErr(c, errors.ErrDBGettingUser, err)
		}
		return
	}

	newPassword := generatePassword(12, 1, 1)
	hashedPassword, err := argon2id.CreateHash(newPassword, argon2id.DefaultParams)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrHashingPassword, err)
		return
	}

	if _, err := h.db.Exec("UPDATE `users` SET `password_hash` = ? WHERE `id` = ?", hashedPassword, userId); err != nil {
		h.services.HandleInternalErr(c, errors.ErrDBUpdatingUserPassword, err)
		return
	}

	if err := h.services.SendNewPassEmail(email, username, newPassword); err != nil {
		h.services.HandleInternalErr(c, errors.ErrSendingEmail, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
