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

const MaxCloakUploadSize = 2 << 10 // 2 kb

func (h *Handler) CloakChangeHandler(c *gin.Context) {
	var input models.CaptchaRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeCloakValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxCloakUploadSize)

	// get cloak
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeCloakValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
			return
		}
	}(file)

	// check cloak size for 22х17 and 64x32
	img, _, err := image.Decode(file)
	if err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if (width != 22 && height != 17) && (width != 64 && height != 32) {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	// check file type
	buffer := make([]byte, fileHeader.Size)
	if _, err = file.Read(buffer); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	fileType := http.DetectContentType(buffer)
	if _, ex := imageTypes[fileType]; !ex {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
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

	if err := h.services.UploadCloak(username, file); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3UploadingCloak, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func (h *Handler) CloakDeleteHandler(c *gin.Context) {
	var input models.CaptchaRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingDeleteCloakValues)
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

	if err := h.services.DeleteCloak(username); err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DeletingCloak, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
