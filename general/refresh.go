package general

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mineleaguedev/rest-api/errors"
	"net/http"
	"strconv"
)

func RefreshHandler(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errors.ErrMissingRefreshToken)
		return
	}
	refreshToken := mapToken["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Middleware.RefreshTokenKey, nil
	})
	if err != nil {
		Service.HandleErr(c, http.StatusUnauthorized, errors.ErrInvalidRefreshToken)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		Service.HandleErr(c, http.StatusUnauthorized, errors.ErrExpiredRefreshToken)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		Service.HandleErr(c, http.StatusUnauthorized, errors.ErrExpiredRefreshToken)
		return
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)
	if !ok {
		Service.HandleErr(c, http.StatusUnprocessableEntity, errors.ErrRefreshTokenUuidNotExists)
		return
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["id"]), 10, 64)
	if err != nil {
		Service.HandleErr(c, http.StatusUnprocessableEntity, errors.ErrRefreshTokenUserIdNotExists)
		return
	}

	td, err := Service.CreateToken(userId)
	if err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrFailedTokenCreation, err)
		return
	}

	if deleted, err := Service.DeleteSession(refreshUuid); err != nil || deleted == 0 {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrDeletingSession, err)
		return
	}

	if err := Service.SaveAuthSession(userId, td); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSavingAuthSession, err)
		return
	}

	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}
