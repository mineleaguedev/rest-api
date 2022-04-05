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
	mapWorldFileName  = "world.rar"
	mapConfigFileName = "map.yml"
	pluginFileName    = "plugin.jar"
	velocityFileName  = "velocity.rar"

	mapWorldFilePath  = "files/" + mapWorldFileName
	mapConfigFilePath = "files/" + mapConfigFileName
	pluginFilePath    = "files/" + pluginFileName
	velocityFilePath  = "files/" + velocityFileName
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
	resp, err := s.config.MiniGamesManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MiniGamesBucket,
		Prefix: aws.String("maps/"),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameMapsList(minigame string) ([]*s3.Object, error) {
	resp, err := s.config.MiniGamesManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MiniGamesBucket,
		Prefix: aws.String("maps/" + minigame),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameFormatMapsList(minigame, format string) ([]*s3.Object, error) {
	resp, err := s.config.MiniGamesManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MiniGamesBucket,
		Prefix: aws.String("maps/" + minigame + "/" + format),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetMiniGameFormatMapVersionsList(minigame, format, mapName string) ([]*s3.Object, error) {
	resp, err := s.config.MiniGamesManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MiniGamesBucket,
		Prefix: aws.String("maps/" + minigame + "/" + format + "/" + mapName),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) UploadMap(minigame, format, mapName, version string, worldFile, configFile multipart.File) error {
	objects := []s3manager.BatchUploadObject{
		{
			Object: &s3manager.UploadInput{
				Bucket: s.config.MiniGamesBucket,
				Key:    aws.String("maps/" + minigame + "/" + format + "/" + mapName + "/" + version + "/" + mapWorldFileName),
				Body:   worldFile,
			},
		},
		{
			Object: &s3manager.UploadInput{
				Bucket: s.config.MiniGamesBucket,
				Key:    aws.String("maps/" + minigame + "/" + format + "/" + mapName + "/" + version + "/" + mapConfigFileName),
				Body:   configFile,
			},
		},
	}

	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	if err := s.config.MiniGamesUploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return err
	}

	return nil
}

func (s *S3Service) DownloadMapWorld(minigame, format, mapName, version string) (*string, *string, error) {
	worldFile, err := os.Create(mapWorldFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer worldFile.Close()

	// world
	_, err = s.config.MiniGamesDownloader.Download(worldFile, &s3.GetObjectInput{
		Bucket: s.config.MiniGamesBucket,
		Key:    aws.String("maps/" + minigame + "/" + format + "/" + mapName + "/" + version + "/" + mapWorldFileName),
	})
	if err != nil {
		return nil, nil, err
	}

	return &mapWorldFilePath, &mapWorldFileName, nil
}

func (s *S3Service) DownloadMapConfig(minigame, format, mapName, version string) (*string, *string, error) {
	mapFile, err := os.Create(mapConfigFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer mapFile.Close()

	// map
	_, err = s.config.MiniGamesDownloader.Download(mapFile, &s3.GetObjectInput{
		Bucket: s.config.MiniGamesBucket,
		Key:    aws.String("maps/" + minigame + "/" + format + "/" + mapName + "/" + version + "/" + mapConfigFileName),
	})
	if err != nil {
		return nil, nil, err
	}

	return &mapConfigFilePath, &mapConfigFileName, nil
}

func (s *S3Service) GetPluginsList() ([]*s3.Object, error) {
	resp, err := s.config.MiniGamesManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MiniGamesBucket,
		Prefix: aws.String("plugins/"),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) GetPluginVersionsList(plugin string) ([]*s3.Object, error) {
	resp, err := s.config.MiniGamesManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MiniGamesBucket,
		Prefix: aws.String("plugins/" + plugin + "/"),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) UploadPlugin(plugin, version string, jarFile multipart.File) error {
	_, err := s.config.MiniGamesUploader.Upload(&s3manager.UploadInput{
		Bucket: s.config.MiniGamesBucket,
		Key:    aws.String("plugins/" + plugin + "/" + version + "/" + pluginFileName),
		Body:   jarFile,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) DownloadPluginJar(plugin, version string) (*string, *string, error) {
	jarFile, err := os.Create(pluginFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer jarFile.Close()

	_, err = s.config.MiniGamesDownloader.Download(jarFile, &s3.GetObjectInput{
		Bucket: s.config.MiniGamesBucket,
		Key:    aws.String("plugins/" + plugin + "/" + version + "/" + pluginFileName),
	})
	if err != nil {
		return nil, nil, err
	}

	return &pluginFilePath, &plugin, nil
}

func (s *S3Service) GetVelocityVersionList() ([]*s3.Object, error) {
	resp, err := s.config.MiniGamesManager.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s.config.MiniGamesBucket,
		Prefix: aws.String("velocity/"),
	})
	if err != nil {
		return nil, err
	}

	return resp.Contents, nil
}

func (s *S3Service) UploadVelocity(version string, rarFile multipart.File) error {
	_, err := s.config.MiniGamesUploader.Upload(&s3manager.UploadInput{
		Bucket: s.config.MiniGamesBucket,
		Key:    aws.String("velocity/" + version + "/" + velocityFileName),
		Body:   rarFile,
	})
	if err != nil {
		return err
	}

	return nil
}
