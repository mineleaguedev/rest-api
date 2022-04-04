package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
)

func (h *Handler) PluginGetHandler(c *gin.Context) {
	plugin := c.Param("name")
	version := c.Param("version")

	filePath, fileName, err := h.services.DownloadPluginJar(plugin, version)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DownloadingPluginJar, err)
		return
	}

	c.FileAttachment(*filePath, *fileName+".jar")
}
