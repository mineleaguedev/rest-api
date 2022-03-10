package general

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/spf13/viper"
	"strings"
)

var (
	regFrom     string
	regSubject  string
	regHtmlBody string
	regTextBody string
	regCharSet  string

	passResetFrom     string
	passResetSubject  string
	passResetHtmlBody string
	passResetTextBody string
	passResetCharSet  string

	emailClient *ses.SES
)

func SetupEmail(client *ses.SES) {
	emailClient = client

	regFrom = viper.GetString("email.reg.from")
	regSubject = viper.GetString("email.reg.subject")
	regHtmlBody = viper.GetString("email.reg.htmlBody")
	regTextBody = viper.GetString("email.reg.textBody")
	regCharSet = viper.GetString("email.reg.charSet")

	passResetFrom = viper.GetString("email.passReset.from")
	passResetSubject = viper.GetString("email.passReset.subject")
	passResetHtmlBody = viper.GetString("email.passReset.htmlBody")
	passResetTextBody = viper.GetString("email.passReset.textBody")
	passResetCharSet = viper.GetString("email.passReset.charSet")
}

func sendRegEmail(to, token string) error {
	htmlBody := strings.Replace(regHtmlBody, "%token%", token, 1)

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
					Charset: aws.String(regCharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(regCharSet),
					Data:    aws.String(regTextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(regCharSet),
				Data:    aws.String(regSubject),
			},
		},
		Source: aws.String(regFrom),
	}

	_, err := emailClient.SendEmail(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return awsErr
		} else {
			return err
		}
	}

	return nil
}

func sendPassResetEmail(to, token string) error {
	htmlBody := strings.Replace(passResetHtmlBody, "%token%", token, 1)

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
					Charset: aws.String(passResetCharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(passResetCharSet),
					Data:    aws.String(passResetTextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(passResetCharSet),
				Data:    aws.String(passResetSubject),
			},
		},
		Source: aws.String(passResetFrom),
	}

	_, err := emailClient.SendEmail(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return awsErr
		} else {
			return err
		}
	}

	return nil
}
