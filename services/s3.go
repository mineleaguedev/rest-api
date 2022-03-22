package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mineleaguedev/rest-api/models"
	"mime/multipart"
	"os"
)

var (
	worldFileName  = "world.rar"
	configFileName = "map.yml"
	pluginFileName = "plugin.jar"

	worldFilePath  = "files/" + worldFileName
	configFilePath = "files/" + configFileName
	pluginFilePath = "files/" + pluginFileName
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
	_, err := s.config.SkinsManager.DeleteObject(&s3.DeleteObjectInput{
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
	_, err := s.config.CloaksManager.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.config.CloaksBucket,
		Key:    aws.String(username + ".png"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) GetMapsList() ([]*s3.Object, error) {
	resp, err := s.config.MapsManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameMapsList(minigame string) ([]*s3.Object, error) {
	resp, err := s.config.MapsManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
		Prefix: aws.String(minigame),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameFormatMapsList(minigame, format string) ([]*s3.Object, error) {
	resp, err := s.config.MapsManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
		Prefix: aws.String(minigame + "/" + format),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameFormatMapVersionsList(minigame, format, mapName string) ([]*s3.Object, error) {
	resp, err := s.config.MapsManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MapsBucket,
		Prefix: aws.String(minigame + "/" + format + "/" + mapName),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) CreateMap(minigame, format, mapName, version string, worldFile, configFile multipart.File) error {
	objects := []s3manager.BatchUploadObject{
		{
			Object: &s3manager.UploadInput{
				Bucket: s.config.MapsBucket,
				Key:    aws.String(minigame + "/" + format + "/" + mapName + "/" + version + "/" + worldFileName),
				Body:   worldFile,
			},
		},
		{
			Object: &s3manager.UploadInput{
				Bucket: s.config.MapsBucket,
				Key:    aws.String(minigame + "/" + format + "/" + mapName + "/" + version + "/" + configFileName),
				Body:   configFile,
			},
		},
	}

	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	if err := s.config.MapsUploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return err
	}

	return nil
}

func (s *S3Service) DownloadMapWorld(minigame, format, mapName, version string) (*string, *string, error) {
	worldFile, err := os.Create(worldFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer worldFile.Close()

	// world
	_, err = s.config.MapsDownloader.Download(worldFile, &s3.GetObjectInput{
		Bucket: s.config.MapsBucket,
		Key:    aws.String(minigame + "/" + format + "/" + mapName + "/" + version + "/" + worldFileName),
	})
	if err != nil {
		return nil, nil, err
	}

	return &worldFilePath, &worldFileName, err
}

func (s *S3Service) DownloadMapConfig(minigame, format, mapName, version string) (*string, *string, error) {
	mapFile, err := os.Create(configFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer mapFile.Close()

	// map
	_, err = s.config.MapsDownloader.Download(mapFile, &s3.GetObjectInput{
		Bucket: s.config.MapsBucket,
		Key:    aws.String(minigame + "/" + format + "/" + mapName + "/" + version + "/" + configFileName),
	})
	if err != nil {
		return nil, nil, err
	}

	return &configFilePath, &configFileName, err
}

func (s *S3Service) GetPluginsList() ([]*s3.Object, error) {
	resp, err := s.config.PluginsManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.PluginsBucket,
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetPluginVersionsList(plugin string) ([]*s3.Object, error) {
	resp, err := s.config.PluginsManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.PluginsBucket,
		Prefix: aws.String(plugin + "/"),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) DownloadPluginJar(plugin, version string) (*string, *string, error) {
	jarFile, err := os.Create(pluginFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer jarFile.Close()

	// world
	_, err = s.config.PluginsDownloader.Download(jarFile, &s3.GetObjectInput{
		Bucket: s.config.PluginsBucket,
		Key:    aws.String(plugin + "/" + version + "/" + pluginFileName),
	})
	if err != nil {
		return nil, nil, err
	}

	return &pluginFilePath, &plugin, err
}
