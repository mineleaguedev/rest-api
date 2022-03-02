package general

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

func MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		middlewareImpl(c)
	}
}

func getClaimsFromJWT(c *gin.Context) (jwt.MapClaims, error) {
	token, err := parseToken(c)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

func middlewareImpl(c *gin.Context) {
	claims, err := getClaimsFromJWT(c)
	if err != nil {
		handleErr(c, http.StatusUnauthorized, err)
		return
	}

	if claims["exp"] == nil {
		handleErr(c, http.StatusBadRequest, ErrMissingExpField)
		return
	}

	if _, ok := claims["exp"].(float64); !ok {
		handleErr(c, http.StatusBadRequest, ErrWrongFormatOfExp)
		return
	}

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		handleErr(c, http.StatusUnauthorized, ErrExpiredToken)
		return
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := identify(c)

	if identity != nil {
		c.Set(Middleware.IdentityKey, identity)
	}

	if !authorize(identity, c) {
		handleErr(c, http.StatusForbidden, ErrForbidden)
		return
	}

	c.Next()
}

func identify(c *gin.Context) interface{} {
	claims := extractClaimsFromContext(c)
	return &User{
		ID:       claims["id"].(float64),
		Username: claims["username"].(string),
	}
}

func authorize(data interface{}, _ *gin.Context) bool {
	if v, ok := data.(*User); ok && v.Username != "" {
		return true
	}
	return false
}
