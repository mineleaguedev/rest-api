package general

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogoutHandler(c *gin.Context) {
	accessDetails, err := extractTokenMetadata(c.Request)
	if err != nil {
		handleErr(c, http.StatusUnauthorized, err)
		return
	}

	if deleted, err := deleteSession(accessDetails.AccessUuid); err != nil || deleted == 0 {
		handleInternalErr(c, http.StatusInternalServerError, ErrDeletingSession, err)
		return
	}

	c.JSON(http.StatusOK, "Successfully logged out")
}
