package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
	"net/http"
)

const MaxVelocityUploadSize = 200 << 20 // 200 mb

type velocityUploadRequest struct {
	Version string `form:"version" binding:"required"`
}

func (h *Handler) VelocityUploadHandler(c *gin.Context) {
	var input velocityUploadRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingVelocityUploadValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxVelocityUploadSize)

	// rar file
	rarFile, rarFileHeader, err := c.Request.FormFile("rarFile")
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingVelocityUploadValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidVelocityRarFile)
			return
		}
	}(rarFile)

	rarFileBuffer := make([]byte, rarFileHeader.Size)
	if _, err = rarFile.Read(rarFileBuffer); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidVelocityRarFile)
		return
	}

	rarFileType := http.DetectContentType(rarFileBuffer)
	fmt.Println(rarFileType)
	if rarFileType != "application/x-rar-compressed" {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidVelocityRarFile)
		return
	}

	if _, err := rarFile.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidVelocityRarFile)
		return
	}

	if err := h.services.UploadVelocity(input.Version, rarFile); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3UploadingVelocity, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
