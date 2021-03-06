package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"regexp"
	"unicode"
)

type registerRequest struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
	Captcha  string `form:"h-captcha-response" binding:"required"`
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

func (h *Handler) RegHandler(c *gin.Context) {
	var input registerRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingRegValues)
		return
	}

	if !validUsername(input.Username) {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidUsername)
		return
	}

	if sevenOrMore, number := validPassword(input.Password); !sevenOrMore || !number {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPassword)
		return
	}

	if response := h.services.VerifyCaptcha(input.Captcha); !response.Success {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	var exists bool
	if err := h.db.QueryRow("SELECT 1 FROM `users` WHERE `username` = ? OR `email` = ?", input.Username, input.Email).Scan(&exists); err != nil && err != sql.ErrNoRows {
		h.services.HandleInternalErr(c, errors.ErrDBRegisteringUser, err)
		return
	}

	if exists {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserAlreadyExists)
		return
	}

	hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrHashingPassword, err)
		return
	}

	regToken := generateToken(40)

	if err := h.services.SaveRegSession(regToken, input.Username, input.Email, hashedPassword, h.middleware.RegTokenTime); err != nil {
		h.services.HandleInternalErr(c, errors.ErrSavingRegSession, err)
		return
	}

	if err := h.services.SendRegEmail(input.Email, regToken); err != nil {
		h.services.HandleInternalErr(c, errors.ErrSendingEmail, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func (h *Handler) RegConfirmHandler(c *gin.Context) {
	token := c.Param("token")

	regInfo, err := h.services.GetRegSession(token)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrGettingRegSession, err)
		return
	}

	if deleted, err := h.services.DeleteSession(token); err != nil || deleted == 0 {
		h.services.HandleInternalErr(c, errors.ErrDeletingSession, err)
		return
	}

	if _, err := h.db.Exec("INSERT INTO `users` (`username`, `email`, `password_hash`) VALUES (?, ?, ?)",
		regInfo.Username, regInfo.Email, regInfo.HashedPassword); err != nil {
		h.services.HandleInternalErr(c, errors.ErrDBRegisteringUser, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
