package general

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func Identify(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		ID:       claims["id"].(float64),
		Username: claims["username"].(string),
	}
}

func Authorize(data interface{}, _ *gin.Context) bool {
	if v, ok := data.(*User); ok && v.Username != "" {
		return true
	}
	return false
}
