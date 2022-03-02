package general

import (
	"github.com/gin-gonic/gin"
	"github.com/kataras/hcaptcha"
	"html/template"
	"log"
)

var (
	SiteKey       string
	CaptchaClient *hcaptcha.Client
	RegForm       *template.Template
	AuthForm      *template.Template
)

func RenderRegForm(c *gin.Context) {
	if err := RegForm.Execute(c.Writer, map[string]string{
		"SiteKey": SiteKey,
	}); err != nil {
		log.Printf("Error rendering reg form: %s\n", err.Error())
		return
	}
}

func RenderAuthForm(c *gin.Context) {
	if err := AuthForm.Execute(c.Writer, map[string]string{
		"SiteKey": SiteKey,
	}); err != nil {
		log.Printf("Error rendering auth form: %s\n", err.Error())
		return
	}
}