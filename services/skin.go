package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
)

type SkinService struct {
	config models.SkinConfig
}

func NewSkinService(skinConfig models.SkinConfig) *SkinService {
	return &SkinService{
		config: skinConfig,
	}
}

func (s *SkinService) SetSkin(username string, file multipart.File) error {
	_, err := s.config.Uploader.Upload(&s3manager.UploadInput{
		Bucket: s.config.Bucket,
		Key:    aws.String(username + ".png"),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SkinService) DeleteSkin(username string) error {
	_, err := s.config.Deleter.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.config.Bucket,
		Key:    aws.String(username + ".png"),
	})
	if err != nil {
		return err
	}

	return nil
}
