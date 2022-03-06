package general

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
)

func RefreshHandler(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrMissingRefreshToken)
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
		handleErr(c, http.StatusUnauthorized, ErrInvalidRefreshToken)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		handleErr(c, http.StatusUnauthorized, ErrExpiredRefreshToken)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		handleErr(c, http.StatusUnauthorized, ErrExpiredRefreshToken)
		return
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)
	if !ok {
		handleErr(c, http.StatusUnprocessableEntity, ErrRefreshTokenUuidNotExists)
		return
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["id"]), 10, 64)
	if err != nil {
		handleErr(c, http.StatusUnprocessableEntity, ErrRefreshTokenUserIdNotExists)
		return
	}

	td, err := createToken(userId)
	if err != nil {
		handleInternalErr(c, http.StatusInternalServerError, ErrFailedTokenCreation, err)
		return
	}

	if deleted, err := deleteSession(refreshUuid); err != nil || deleted == 0 {
		handleInternalErr(c, http.StatusInternalServerError, ErrDeletingTokenSession, err)
		return
	}

	if err := saveSession(userId, td); err != nil {
		handleInternalErr(c, http.StatusInternalServerError, ErrSavingTokenSession, err)
		return
	}

	tokens := map[string]string{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}
