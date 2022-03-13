package general

import (
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
)

func ChangePassHandler(c *gin.Context) {
	var input models.ChangePassRequest

	if err := c.ShouldBind(&input); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangePassValues)
		return
	}

	if sevenOrMore, number := validPassword(input.NewPassword); !sevenOrMore || !number {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPassword)
		return
	}

	if response := Service.VerifyCaptcha(input.Captcha); !response.Success {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var email string
	var oldHashedPassword string
	err := DB.QueryRow("SELECT `email`, `password_hash` FROM `users` WHERE `id` = ?", userId).Scan(&email, &oldHashedPassword)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist, err)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.OldPassword, oldHashedPassword)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrUnhashingPassword, err)
		return
	}

	if !match {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrWrongPassword)
		return
	}

	newHashedPassword, err := argon2id.CreateHash(input.NewPassword, argon2id.DefaultParams)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrHashingPassword, err)
		return
	}

	if _, err := DB.Exec("UPDATE `users` SET `password_hash` = ? WHERE `id` = ?", newHashedPassword, userId); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrUpdatingUserPassword, err)
		return
	}

	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	if err := Service.SendChangePassEmail(email, ip); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSendingEmail, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
