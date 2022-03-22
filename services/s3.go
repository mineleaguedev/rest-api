package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
)

type S3Service struct {
	config models.S3Config
}

func NewS3Service(s3Config models.S3Config) *S3Service {
	return &S3Service{
		config: s3Config,
	}
}

func (s *S3Service) UploadSkin(username string, file multipart.File) error {
	_, err := s.config.SkinsUploader.Upload(&s3manager.UploadInput{
		Bucket: s.config.SkinsBucket,
		Key:    aws.String(username + ".png"),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) DeleteSkin(username string) error {
	_, err := s.config.SkinsDeleter.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.config.SkinsBucket,
		Key:    aws.String(username + ".png"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) UploadCloak(username string, file multipart.File) error {
	_, err := s.config.CloaksUploader.Upload(&s3manager.UploadInput{
		Bucket: s.config.CloaksBucket,
		Key:    aws.String(username + ".png"),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) DeleteCloak(username string) error {
	_, err := s.config.CloaksDeleter.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.config.CloaksBucket,
		Key:    aws.String(username + ".png"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) GetMapsList() ([]*s3.Object, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
	}
	resp, err := s.config.MapsDeleter.ListObjectsV2(params)
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameMapsList(minigame string) ([]*s3.Object, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
		Prefix: aws.String(minigame),
	}
	resp, err := s.config.MapsDeleter.ListObjectsV2(params)
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameFormatMapsList(minigame, format string) ([]*s3.Object, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
		Prefix: aws.String(minigame + "/" + format),
	}
	resp, err := s.config.MapsDeleter.ListObjectsV2(params)
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameFormatMapVersionsList(minigame, format, mapName string) ([]*s3.Object, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
		Prefix: aws.String(minigame + "/" + format + "/" + mapName),
	}
	resp, err := s.config.MapsDeleter.ListObjectsV2(params)
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}
