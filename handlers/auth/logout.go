package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
)

func (h *Handler) LogoutHandler(c *gin.Context) {
	accessDetails, err := h.services.ExtractTokenMetadata(c.Request)
	if err != nil {
		h.services.HandleErr(c, http.StatusUnauthorized, err)
		return
	}

	if deleted, err := h.services.DeleteSession(accessDetails.AccessUuid); err != nil || deleted == 0 {
		h.services.HandleInternalErr(c, errors.ErrDeletingSession, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
