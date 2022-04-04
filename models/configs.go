package models

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
	goredis "github.com/go-redis/redis/v8"
	"github.com/kataras/hcaptcha"
	"github.com/nitishm/go-rejson/v4"
	"html/template"
)

type RedisConfig struct {
	Client      *goredis.Client
	JsonHandler *rejson.Handler
	Ctx         context.Context
}

type EmailConfig struct {
	RegFrom            string
	RegSubject         string
	RegHtmlBody        string
	PassResetFrom      string
	PassResetSubject   string
	PassResetHtmlBody  string
	NewPassFrom        string
	NewPassSubject     string
	NewPassHtmlBody    string
	ChangePassFrom     string
	ChangePassSubject  string
	ChangePassHtmlBody string
	Client             *ses.SES
}

type CaptchaConfig struct {
	SiteKey         string
	Client          *hcaptcha.Client
	RegForm         *template.Template
	AuthForm        *template.Template
	PassResetForm   *template.Template
	ChangePassForm  *template.Template
	ChangeSkinForm  *template.Template
	DeleteSkinForm  *template.Template
	ChangeCloakForm *template.Template
	DeleteCloakForm *template.Template
}

type S3Config struct {
	SkinsBucket         *string
	SkinsUploader       *s3manager.Uploader
	SkinsManager        *s3.S3
	CloaksBucket        *string
	CloaksUploader      *s3manager.Uploader
	CloaksManager       *s3.S3
	MiniGamesBucket     *string
	MiniGamesUploader   *s3manager.Uploader
	MiniGamesDownloader *s3manager.Downloader
	MiniGamesManager    *s3.S3
}
