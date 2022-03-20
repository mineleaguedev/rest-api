package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"net/http"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessDetails, err := h.services.ExtractTokenMetadata(c.Request)
		if err != nil {
			h.services.HandleErr(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		userId, err := h.services.GetAuthSession(accessDetails)
		if err != nil {
			h.services.HandleInternalErr(c, errors.ErrGettingAuthSession, err)
			c.Abort()
			return
		}

		c.Set("userId", userId)

		c.Next()
	}
}
