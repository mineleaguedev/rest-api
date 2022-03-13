package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mineleaguedev/rest-api/models"
	"strings"
)

type EmailService struct {
	config models.EmailConfig
}

func NewEmailService(emailConfig models.EmailConfig) *EmailService {
	return &EmailService{
		config: emailConfig,
	}
}

func (s *EmailService) SendRegEmail(to, token string) error {
	htmlBody := strings.Replace(s.config.RegHtmlBody, "%token%", token, 1)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(s.config.RegCharSet),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(s.config.RegCharSet),
				Data:    aws.String(s.config.RegSubject),
			},
		},
		Source: aws.String(s.config.RegFrom),
	}

	_, err := s.config.Client.SendEmail(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return awsErr
		} else {
			return err
		}
	}

	return nil
}

func (s *EmailService) SendPassResetEmail(to, token, username, ipAddress string) error {
	replacer := strings.NewReplacer("%pass_reset_token%", token, "%username%", username, "%ip_address%", ipAddress)
	htmlBody := replacer.Replace(s.config.PassResetHtmlBody)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(s.config.PassResetCharSet),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(s.config.PassResetCharSet),
				Data:    aws.String(s.config.PassResetSubject),
			},
		},
		Source: aws.String(s.config.PassResetFrom),
	}

	_, err := s.config.Client.SendEmail(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return awsErr
		} else {
			return err
		}
	}

	return nil
}

func (s *EmailService) SendNewPassEmail(to, username, password string) error {
	replacer := strings.NewReplacer("%username%", username, "%password%", password)
	htmlBody := replacer.Replace(s.config.NewPassHtmlBody)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(s.config.NewPassCharSet),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(s.config.NewPassCharSet),
				Data:    aws.String(s.config.NewPassSubject),
			},
		},
		Source: aws.String(s.config.NewPassFrom),
	}

	_, err := s.config.Client.SendEmail(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return awsErr
		} else {
			return err
		}
	}

	return nil
}
