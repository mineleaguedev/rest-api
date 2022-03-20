package models

type PassChangeRequest struct {
	OldPassword string `form:"old-password" binding:"required"`
	NewPassword string `form:"new-password" binding:"required"`
	Captcha     string `form:"h-captcha-response" binding:"required"`
}

type MoneyTransferRequest struct {
	Username string `form:"username" binding:"required"`
	Amount   int64  `form:"amount" binding:"required"`
	Captcha  string `form:"h-captcha-response" binding:"required"`
}

type CaptchaRequest struct {
	Captcha string `form:"h-captcha-response" binding:"required"`
}
