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

func (s *EmailService) SendPassResetEmail(to, token string) error {
	htmlBody := strings.Replace(s.config.PassResetHtmlBody, "%token%", token, 1)

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
