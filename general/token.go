package general

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/twinj/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func createToken(userId int64) (*TokenDetails, error) {
	td := TokenDetails{}
	td.AtExpires = time.Now().Add(Middleware.AccessTokenTime).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(Middleware.RefreshTokenTime).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{
		"access_uuid": td.AccessUuid,
		"id":          userId,
		"exp":         td.AtExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(Middleware.AccessTokenKey)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{
		"refresh_uuid": td.RefreshUuid,
		"id":           userId,
		"exp":          td.RtExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(Middleware.RefreshTokenKey)
	if err != nil {
		return nil, err
	}

	return &td, nil
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Middleware.AccessTokenKey, nil
	})
	if err != nil {
		return nil, ErrInvalidAccessToken
	}

	return token, nil
}

func extractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrExpiredAccessToken
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, ErrAccessTokenUuidNotExists
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["id"]), 10, 64)
	if err != nil {
		return nil, ErrAccessTokenUserIdNotExists
	}

	return &AccessDetails{
		AccessUuid: accessUuid,
		UserId:     userId,
	}, nil
}
