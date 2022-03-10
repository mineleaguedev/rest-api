package general

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func AuthHandler(c *gin.Context) {
	userId, httpCode, err := authenticate(c)
	if err != nil {
		handleErr(c, httpCode, err)
		return
	}

	td, err := createToken(*userId)
	if err != nil {
		handleInternalErr(c, http.StatusInternalServerError, ErrFailedTokenCreation, err)
		return
	}

	if err := saveAuthSession(*userId, td); err != nil {
		handleInternalErr(c, http.StatusInternalServerError, ErrSavingAuthSession, err)
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
		return nil, http.StatusBadRequest, ErrMissingAuthValues
	}

	if response := captchaClient.VerifyToken(input.Captcha); !response.Success {
		return nil, http.StatusBadRequest, ErrInvalidCaptcha
	}

	var id int64
	var hashedPassword string
	err := DB.QueryRow("SELECT `id`, `password_hash` FROM `users` WHERE `username` = ?", input.Username).Scan(&id, &hashedPassword)
	if err != nil && err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, ErrUserDoesNotExist
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, hashedPassword)
	if err != nil {
		log.Printf(ErrUnhashingPassword.Error()+": %s\n", err.Error())
		return nil, http.StatusInternalServerError, ErrUnhashingPassword
	}

	if !match {
		return nil, http.StatusBadRequest, ErrWrongUsernameOrPassword
	}

	return &id, http.StatusOK, nil
}
