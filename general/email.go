package general

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
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
}

func sendRegEmail(to string) error {
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
					Data:    aws.String(regHtmlBody),
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

func sendPassResetEmail(to string) error {
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
					Data:    aws.String(passResetHtmlBody),
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
