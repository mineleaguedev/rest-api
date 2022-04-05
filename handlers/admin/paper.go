package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
	"net/http"
)

const MaxPaperUploadSize = 200 << 20 // 200 mb

type paperUploadRequest struct {
	Version string `form:"version" binding:"required"`
}

func (h *Handler) PaperUploadHandler(c *gin.Context) {
	var input paperUploadRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPaperUploadValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxPaperUploadSize)

	// rar file
	rarFile, rarFileHeader, err := c.Request.FormFile("rarFile")
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingPaperUploadValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPaperRarFile)
			return
		}
	}(rarFile)

	rarFileBuffer := make([]byte, rarFileHeader.Size)
	if _, err = rarFile.Read(rarFileBuffer); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPaperRarFile)
		return
	}

	rarFileType := http.DetectContentType(rarFileBuffer)
	if rarFileType != "application/x-rar-compressed" {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPaperRarFile)
		return
	}

	if _, err := rarFile.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidPaperRarFile)
		return
	}

	if err := h.services.UploadPaper(input.Version, rarFile); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3UploadingPaper, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
