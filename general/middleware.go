package general

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessDetails, err := extractTokenMetadata(c.Request)
		if err != nil {
			handleErr(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		userId, err := getSession(accessDetails)
		if err != nil {
			handleInternalErr(c, http.StatusInternalServerError, ErrGettingTokenSession, err)
			c.Abort()
			return
		}

		c.Set("userId", userId)

		c.Next()
	}
}
