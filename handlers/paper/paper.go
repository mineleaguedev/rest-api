package paper

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
)

func (h *Handler) PaperGetHandler(c *gin.Context) {
	version := c.Param("version")

	filePath, fileName, err := h.services.DownloadPaper(version)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DownloadingPaperRar, err)
		return
	}

	c.FileAttachment(*filePath, *fileName)
}
