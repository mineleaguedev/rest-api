package general

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
)

func AuthHandler(c *gin.Context) {
	userId, httpCode, err := authenticate(c)
	if err != nil {
		Service.HandleErr(c, httpCode, err)
		return
	}

	td, err := Service.CreateToken(*userId)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrFailedTokenCreation, err)
		return
	}

	if err := Service.SaveAuthSession(*userId, td); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSavingAuthSession, err)
		return
	}

	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func authenticate(c *gin.Context) (*int64, int, error) {
	var input models.AuthRequest

	if err := c.ShouldBind(&input); err != nil {
		return nil, http.StatusBadRequest, errors.ErrMissingAuthValues
	}

	if response := Service.VerifyCaptcha(input.Captcha); !response.Success {
		return nil, http.StatusBadRequest, errors.ErrInvalidCaptcha
	}

	var id int64
	var hashedPassword string
	err := DB.QueryRow("SELECT `id`, `password_hash` FROM `users` WHERE `username` = ?", input.Username).Scan(&id, &hashedPassword)
	if err != nil && err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, errors.ErrUserDoesNotExist
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, hashedPassword)
	if err != nil {
		log.Printf(errors.ErrUnhashingPassword.Error()+": %s\n", err.Error())
		return nil, http.StatusInternalServerError, errors.ErrUnhashingPassword
	}

	if !match {
		return nil, http.StatusBadRequest, errors.ErrWrongUsernameOrPassword
	}

	return &id, http.StatusOK, nil
}
