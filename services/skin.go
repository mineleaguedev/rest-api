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

func (s *SkinService) UploadSkin(username string, file multipart.File) error {
	_, err := s.config.SkinUploader.Upload(&s3manager.UploadInput{
		Bucket: s.config.SkinBucket,
		Key:    aws.String(username + ".png"),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SkinService) DeleteSkin(username string) error {
	_, err := s.config.SkinDeleter.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.config.SkinBucket,
		Key:    aws.String(username + ".png"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SkinService) UploadCloak(username string, file multipart.File) error {
	_, err := s.config.CloakUploader.Upload(&s3manager.UploadInput{
		Bucket: s.config.CloakBucket,
		Key:    aws.String(username + ".png"),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SkinService) DeleteCloak(username string) error {
	_, err := s.config.CloakDeleter.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.config.CloakBucket,
		Key:    aws.String(username + ".png"),
	})
	if err != nil {
		return err
	}

	return nil
}
