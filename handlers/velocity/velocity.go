package velocity

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
)

func (h *Handler) VelocityGetHandler(c *gin.Context) {
	version := c.Param("version")

	filePath, fileName, err := h.services.DownloadVelocity(version)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DownloadingVelocityRar, err)
		return
	}

	c.FileAttachment(*filePath, *fileName)
}
