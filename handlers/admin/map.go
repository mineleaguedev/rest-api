package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
	"net/http"
)

const MaxMapUploadSize = 200 << 20 // 200 mb

type mapCreateRequest struct {
	MiniGame string `form:"minigame" binding:"required"`
	Format   string `form:"format" binding:"required"`
	Map      string `form:"map" binding:"required"`
	Version  string `form:"version" binding:"required"`
}

func (h *Handler) MapUploadHandler(c *gin.Context) {
	var input mapCreateRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingMapUploadValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxMapUploadSize)

	// world file
	worldFile, worldFileHeader, err := c.Request.FormFile("worldFile")
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingMapUploadValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapWorldFile)
			return
		}
	}(worldFile)

	worldFileBuffer := make([]byte, worldFileHeader.Size)
	if _, err = worldFile.Read(worldFileBuffer); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapWorldFile)
		return
	}

	worldFileType := http.DetectContentType(worldFileBuffer)
	if worldFileType != "application/x-rar-compressed" {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapWorldFile)
		return
	}

	if _, err := worldFile.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapWorldFile)
		return
	}

	// config file
	configFile, configFileHeader, err := c.Request.FormFile("configFile")
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingMapUploadValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapConfigFile)
			return
		}
	}(configFile)

	configFileBuffer := make([]byte, configFileHeader.Size)
	if _, err = configFile.Read(configFileBuffer); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapConfigFile)
		return
	}

	configFileType := http.DetectContentType(configFileBuffer)
	if configFileType != "text/plain; charset=utf-8" {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapConfigFile)
		return
	}

	if _, err := configFile.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidMapConfigFile)
		return
	}

	if err := h.services.UploadMap(input.MiniGame, input.Format, input.Map, input.Version, worldFile, configFile); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3UploadingMap, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
