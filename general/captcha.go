package general

import (
	"github.com/gin-gonic/gin"
	"github.com/kataras/hcaptcha"
	"html/template"
	"log"
)

var (
	captchaSiteKey string
	captchaClient  *hcaptcha.Client
	regForm        *template.Template
	authForm       *template.Template
)

func SetupCaptcha(siteKey, secretKey string) {
	captchaSiteKey = siteKey
	captchaClient = hcaptcha.New(secretKey)
	regForm = template.Must(template.ParseFiles("./forms/reg_form.html"))
	authForm = template.Must(template.ParseFiles("./forms/auth_form.html"))
}

func RenderRegForm(c *gin.Context) {
	if err := regForm.Execute(c.Writer, map[string]string{
		"SiteKey": captchaSiteKey,
	}); err != nil {
		log.Printf("Error rendering reg form: %s\n", err.Error())
		return
	}
}

func RenderAuthForm(c *gin.Context) {
	if err := authForm.Execute(c.Writer, map[string]string{
		"SiteKey": captchaSiteKey,
	}); err != nil {
		log.Printf("Error rendering auth form: %s\n", err.Error())
		return
	}
}
