package services

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/hcaptcha"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
	"net/http"
	"time"
)

type Token interface {
	CreateToken(userId int64) (*models.TokenDetails, error)
	ExtractToken(r *http.Request) string
	VerifyToken(r *http.Request) (*jwt.Token, error)
	ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error)
}

type Err interface {
	HandleErr(c *gin.Context, httpCode int, err error)
	HandleInternalErr(c *gin.Context, httpCode int, err, internalErr error)
}

type Redis interface {
	SaveRegSession(token, username, email, hashedPassword string, expireTime int64) error
	GetRegSession(token string) (*regInfo, error)
	SavePassResetSession(token string, userId int64, expireTime time.Duration) error
	GetPassResetSession(token string) (int64, error)
	SaveAuthSession(userId int64, td *models.TokenDetails) error
	GetAuthSession(accessDetails *models.AccessDetails) (int64, error)
	DeleteSession(key string) (int64, error)
}

type Email interface {
	SendRegEmail(to, token string) error
	SendPassResetEmail(to, token, username, ipAddress string) error
	SendNewPassEmail(to, username, password string) error
	SendChangePassEmail(to, ipAddress string) error
}

type Captcha interface {
	RenderRegForm(c *gin.Context)
	RenderAuthForm(c *gin.Context)
	RenderPassResetForm(c *gin.Context)
	RenderChangePassForm(c *gin.Context)
	RenderChangeSkinForm(c *gin.Context)
	RenderDeleteSkinForm(c *gin.Context)
	RenderChangeCloakForm(c *gin.Context)
	RenderDeleteCloakForm(c *gin.Context)
	VerifyCaptcha(token string) (response hcaptcha.Response)
}

type Skin interface {
	SetSkin(username string, file multipart.File) error
	DeleteSkin(username string) error
	SetCloak(username string, file multipart.File) error
	DeleteCloak(username string) error
}

type Service struct {
	Token
	Err
	Redis
	Email
	Captcha
	Skin
}

func NewService(jwtMiddleware models.JWTMiddleware, redisConfig models.RedisConfig, emailConfig models.EmailConfig,
	captchaConfig models.CaptchaConfig, skinConfig models.SkinConfig) *Service {
	return &Service{
		Token:   NewTokenService(jwtMiddleware),
		Err:     NewErrService(),
		Redis:   NewRedisService(redisConfig),
		Email:   NewEmailService(emailConfig),
		Captcha: NewCaptchaService(captchaConfig),
		Skin:    NewSkinService(skinConfig),
	}
}
