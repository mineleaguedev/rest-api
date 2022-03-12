package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"github.com/twinj/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TokenService struct {
	middleware models.JWTMiddleware
}

func NewTokenService(jwtMiddleware models.JWTMiddleware) *TokenService {
	return &TokenService{
		middleware: jwtMiddleware,
	}
}

func (s *TokenService) CreateToken(userId int64) (*models.TokenDetails, error) {
	td := models.TokenDetails{}
	td.AtExpires = time.Now().Add(s.middleware.AccessTokenTime).Unix()
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(s.middleware.RefreshTokenTime).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{
		"access_uuid": td.AccessUuid,
		"id":          userId,
		"exp":         td.AtExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(s.middleware.AccessTokenKey)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{
		"refresh_uuid": td.RefreshUuid,
		"id":           userId,
		"exp":          td.RtExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(s.middleware.RefreshTokenKey)
	if err != nil {
		return nil, err
	}

	return &td, nil
}

func (s *TokenService) ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (s *TokenService) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := s.ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.middleware.AccessTokenKey, nil
	})
	if err != nil {
		return nil, errors.ErrInvalidAccessToken
	}

	return token, nil
}

func (s *TokenService) ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := s.VerifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrExpiredAccessToken
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, errors.ErrAccessTokenUuidNotExists
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["id"]), 10, 64)
	if err != nil {
		return nil, errors.ErrAccessTokenUserIdNotExists
	}

	return &models.AccessDetails{
		AccessUuid: accessUuid,
		UserId:     userId,
	}, nil
}
