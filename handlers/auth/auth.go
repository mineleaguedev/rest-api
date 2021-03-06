package auth

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"net/http"
)

type authRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Captcha  string `form:"h-captcha-response" binding:"required"`
}

func (h *Handler) AuthHandler(c *gin.Context) {
	var input authRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingAuthValues)
		return
	}

	if response := h.services.VerifyCaptcha(input.Captcha); !response.Success {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	var userId int64
	var hashedPassword string
	if err := h.db.QueryRow("SELECT `id`, `password_hash` FROM `users` WHERE `username` = ?", input.Username).Scan(&userId, &hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleInternalErr(c, errors.ErrDBGettingUser, err)
		}
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, hashedPassword)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrUnhashingPassword, err)
		return
	}

	if !match {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrWrongUsernameOrPassword)
		return
	}

	td, err := h.services.CreateToken(userId)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrFailedTokenCreation, err)
		return
	}

	if err := h.services.SaveAuthSession(userId, td); err != nil {
		h.services.HandleInternalErr(c, errors.ErrSavingAuthSession, err)
		return
	}

	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}
