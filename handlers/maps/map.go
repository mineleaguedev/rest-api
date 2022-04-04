package maps

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
)

func (h *Handler) MapWorldGetHandler(c *gin.Context) {
	minigame := c.Param("minigame")
	format := c.Param("format")
	mapName := c.Param("map")
	version := c.Param("version")

	filePath, fileName, err := h.services.DownloadMapWorld(minigame, format, mapName, version)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DownloadingMapWorld, err)
		return
	}

	c.FileAttachment(*filePath, *fileName)
}

func (h *Handler) MapConfigGetHandler(c *gin.Context) {
	minigame := c.Param("minigame")
	format := c.Param("format")
	mapName := c.Param("map")
	version := c.Param("version")

	filePath, fileName, err := h.services.DownloadMapConfig(minigame, format, mapName, version)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DownloadingMapConfig, err)
		return
	}

	c.FileAttachment(*filePath, *fileName)
}
