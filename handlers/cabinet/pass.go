package cabinet

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"unicode"
)

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

func (h *Handler) ChangePassHandler(c *gin.Context) {
	var input models.ChangePassRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangePassValues)
		return
	}

	if sevenOrMore, number := validPassword(input.NewPassword); !sevenOrMore || !number {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPassword)
		return
	}

	if response := h.services.VerifyCaptcha(input.Captcha); !response.Success {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var email, oldHashedPassword string
	if err := h.db.QueryRow("SELECT `email`, `password_hash` FROM `users` WHERE `id` = ?", userId).Scan(&email, &oldHashedPassword); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleDBErr(c, err)
		}
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.OldPassword, oldHashedPassword)
	if err != nil {
		h.services.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrUnhashingPassword, err)
		return
	}

	if !match {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrWrongPassword)
		return
	}

	newHashedPassword, err := argon2id.CreateHash(input.NewPassword, argon2id.DefaultParams)
	if err != nil {
		h.services.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrHashingPassword, err)
		return
	}

	if _, err := h.db.Exec("UPDATE `users` SET `password_hash` = ? WHERE `id` = ?", newHashedPassword, userId); err != nil {
		h.services.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrDBUpdatingUserPassword, err)
		return
	}

	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	if err := h.services.SendChangePassEmail(email, ip); err != nil {
		h.services.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSendingEmail, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
