package services

import (
	"github.com/gin-gonic/gin"
	"github.com/kataras/hcaptcha"
	"github.com/mineleaguedev/rest-api/models"
	"log"
)

type CaptchaService struct {
	config models.CaptchaConfig
}

func NewCaptchaService(captchaConfig models.CaptchaConfig) *CaptchaService {
	return &CaptchaService{
		config: captchaConfig,
	}
}

func (s *CaptchaService) RenderRegForm(c *gin.Context) {
	if err := s.config.RegForm.Execute(c.Writer, map[string]string{
		"SiteKey": s.config.SiteKey,
	}); err != nil {
		log.Printf("Error rendering reg form: %s\n", err.Error())
		return
	}
}

func (s *CaptchaService) RenderAuthForm(c *gin.Context) {
	if err := s.config.AuthForm.Execute(c.Writer, map[string]string{
		"SiteKey": s.config.SiteKey,
	}); err != nil {
		log.Printf("Error rendering auth form: %s\n", err.Error())
		return
	}
}

func (s *CaptchaService) RenderPassResetForm(c *gin.Context) {
	if err := s.config.PassResetForm.Execute(c.Writer, map[string]string{
		"SiteKey": s.config.SiteKey,
	}); err != nil {
		log.Printf("Error rendering password reset form: %s\n", err.Error())
		return
	}
}

func (s *CaptchaService) RenderChangePassForm(c *gin.Context) {
	if err := s.config.ChangePassForm.Execute(c.Writer, map[string]string{
		"SiteKey": s.config.SiteKey,
	}); err != nil {
		log.Printf("Error rendering change password form: %s\n", err.Error())
		return
	}
}

func (s *CaptchaService) RenderChangeSkinForm(c *gin.Context) {
	if err := s.config.ChangeSkinForm.Execute(c.Writer, map[string]string{
		"SiteKey": s.config.SiteKey,
	}); err != nil {
		log.Printf("Error rendering change skin form: %s\n", err.Error())
		return
	}
}

func (s *CaptchaService) RenderDeleteSkinForm(c *gin.Context) {
	if err := s.config.DeleteSkinForm.Execute(c.Writer, map[string]string{
		"SiteKey": s.config.SiteKey,
	}); err != nil {
		log.Printf("Error rendering delete skin form: %s\n", err.Error())
		return
	}
}

func (s *CaptchaService) VerifyCaptcha(token string) (response hcaptcha.Response) {
	return s.config.Client.VerifyToken(token)
}
