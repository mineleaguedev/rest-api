package general

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
)

const (
	MaxCloakUploadSize = 2 << 10 // 2 kb
)

func ChangeCloakHandler(c *gin.Context) {
	var input models.ChangeSkinRequest

	if err := c.ShouldBind(&input); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeCloakValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxCloakUploadSize)

	// get cloak
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeCloakValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
			return
		}
	}(file)

	// check cloak size for 22х17 and 64x32
	img, _, err := image.Decode(file)
	if err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if (width != 22 && height != 17) && (width != 64 && height != 32) {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	// check file type
	buffer := make([]byte, fileHeader.Size)
	if _, err = file.Read(buffer); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	fileType := http.DetectContentType(buffer)
	if _, ex := ImageTypes[fileType]; !ex {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCloak)
		return
	}

	// other checks
	if response := Service.VerifyCaptcha(input.Captcha); !response.Success {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var username string
	if err := DB.QueryRow("SELECT `username` FROM `users` WHERE `id` = ?", userId).Scan(&username); err != nil {
		Service.HandleInternalErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist, err)
		return
	}

	if err := Service.SetCloak(username, file); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSettingCloak, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func DeleteCloakHandler(c *gin.Context) {
	var input models.ChangeSkinRequest

	if err := c.ShouldBind(&input); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeCloakValues)
		return
	}

	if response := Service.VerifyCaptcha(input.Captcha); !response.Success {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var username string
	if err := DB.QueryRow("SELECT `username` FROM `users` WHERE `id` = ?", userId).Scan(&username); err != nil {
		Service.HandleInternalErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist, err)
		return
	}

	if err := Service.DeleteCloak(username); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrDeletingCloak, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
