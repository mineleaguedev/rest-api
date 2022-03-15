package models

type RegisterRequest struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
	Captcha  string `form:"h-captcha-response" binding:"required"`
}

type AuthRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Captcha  string `form:"h-captcha-response" binding:"required"`
}

type PassResetRequest struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Captcha  string `form:"h-captcha-response" binding:"required"`
}

type ChangePassRequest struct {
	OldPassword string `form:"old-password" binding:"required"`
	NewPassword string `form:"new-password" binding:"required"`
	Captcha     string `form:"h-captcha-response" binding:"required"`
}

type ChangeSkinRequest struct {
	Captcha string `form:"h-captcha-response" binding:"required"`
}
