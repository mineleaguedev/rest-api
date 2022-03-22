package cabinet

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
)

const MaxSkinUploadSize = 4 << 10 // 4 kb

func (h *Handler) SkinChangeHandler(c *gin.Context) {
	var input models.CaptchaRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeSkinValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxSkinUploadSize)

	// get skin
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeSkinValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
			return
		}
	}(file)

	// check skin size for 64x64 and 64x32
	img, _, err := image.Decode(file)
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width != 64 || (height != 64 && height != 32) {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	// check file type
	buffer := make([]byte, fileHeader.Size)
	if _, err = file.Read(buffer); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	fileType := http.DetectContentType(buffer)
	if _, ex := imageTypes[fileType]; !ex {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	// other checks
	if response := h.services.VerifyCaptcha(input.Captcha); !response.Success {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var username string
	if err := h.db.QueryRow("SELECT `username` FROM `users` WHERE `id` = ?", userId).Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleInternalErr(c, errors.ErrDBGettingUser, err)
		}
		return
	}

	if err := h.services.UploadSkin(username, file); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3UploadingSkin, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func (h *Handler) SkinDeleteHandler(c *gin.Context) {
	var input models.CaptchaRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingDeleteSkinValues)
		return
	}

	if response := h.services.VerifyCaptcha(input.Captcha); !response.Success {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var username string
	if err := h.db.QueryRow("SELECT `username` FROM `users` WHERE `id` = ?", userId).Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleInternalErr(c, errors.ErrDBGettingUser, err)
		}
		return
	}

	if err := h.services.DeleteSkin(username); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DeletingSkin, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
