package models

type CaptchaRequest struct {
	Captcha string `form:"h-captcha-response" binding:"required"`
}
