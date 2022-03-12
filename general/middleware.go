package general

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessDetails, err := Service.ExtractTokenMetadata(c.Request)
		if err != nil {
			Service.HandleErr(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		userId, err := Service.GetAuthSession(accessDetails)
		if err != nil {
			Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrGettingAuthSession, err)
			c.Abort()
			return
		}

		c.Set("userId", userId)

		c.Next()
	}
}
