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
	MaxUploadSize = 4 << 10 // 4 kb
)

var (
	ImageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
	}
)

func ChangeSkinHandler(c *gin.Context) {
	var input models.ChangeSkinRequest

	if err := c.ShouldBind(&input); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeSkinValues)
		return
	}

	// limit upload file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	// get skin
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeSkinValues)
		return
	}

	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
			return
		}
	}(file)

	// check skin size for 64x64 and 64x32
	img, _, err := image.Decode(file)
	if err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width != 64 || (height != 64 && height != 32) {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	// check file type
	buffer := make([]byte, fileHeader.Size)
	if _, err = file.Read(buffer); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	fileType := http.DetectContentType(buffer)
	if _, ex := ImageTypes[fileType]; !ex {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
		return
	}

	// seek file
	if _, err := file.Seek(0, 0); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidSkin)
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

	if err := Service.SetSkin(username, file); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSettingSkin, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}

func DeleteSkinHandler(c *gin.Context) {
	var input models.ChangeSkinRequest

	if err := c.ShouldBind(&input); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingChangeSkinValues)
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

	if err := Service.DeleteSkin(username); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrDeletingSkin, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
