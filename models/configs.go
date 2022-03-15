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
	RegCharSet         string
	PassResetFrom      string
	PassResetSubject   string
	PassResetHtmlBody  string
	PassResetCharSet   string
	NewPassFrom        string
	NewPassSubject     string
	NewPassHtmlBody    string
	NewPassCharSet     string
	ChangePassFrom     string
	ChangePassSubject  string
	ChangePassHtmlBody string
	ChangePassCharSet  string
	Client             *ses.SES
}

type CaptchaConfig struct {
	SiteKey        string
	Client         *hcaptcha.Client
	RegForm        *template.Template
	AuthForm       *template.Template
	PassResetForm  *template.Template
	ChangePassForm *template.Template
	ChangeSkinForm *template.Template
	DeleteSkinForm *template.Template
}

type SkinConfig struct {
	Bucket   *string
	Uploader *s3manager.Uploader
	Deleter  *s3.S3
}
