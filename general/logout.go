package general

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"net/http"
)

func LogoutHandler(c *gin.Context) {
	accessDetails, err := Service.ExtractTokenMetadata(c.Request)
	if err != nil {
		Service.HandleErr(c, http.StatusUnauthorized, err)
		return
	}

	if deleted, err := Service.DeleteSession(accessDetails.AccessUuid); err != nil || deleted == 0 {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrDeletingSession, err)
		return
	}

	c.JSON(http.StatusOK, "Successfully logged out")
}
