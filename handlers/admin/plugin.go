package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
	"net/http"
)

const MaxPluginUploadSize = 200 << 20 // 200 mb

type pluginUploadRequest struct {
	Plugin  string `form:"plugin" binding:"required"`
	Version string `form:"version" binding:"required"`
}

func (h *Handler) PluginUploadHandler(c *gin.Context) {
	var input pluginUploadRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPluginUploadValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxPluginUploadSize)

	// jar file
	jarFile, jarFileHeader, err := c.Request.FormFile("jarFile")
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPluginUploadValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPluginJarFile)
			return
		}
	}(jarFile)

	jarFileBuffer := make([]byte, jarFileHeader.Size)
	if _, err = jarFile.Read(jarFileBuffer); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPluginJarFile)
		return
	}

	jarFileType := http.DetectContentType(jarFileBuffer)
	if jarFileType != "application/zip" {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPluginJarFile)
		return
	}

	if _, err := jarFile.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPluginJarFile)
		return
	}

	if err := h.services.UploadPlugin(input.Plugin, input.Version, jarFile); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3UploadingPlugin, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
