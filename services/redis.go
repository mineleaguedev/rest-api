package services

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/mineleaguedev/rest-api/models"
	"strconv"
	"time"
)

type regInfo struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
}

type RedisService struct {
	config models.RedisConfig
}

func NewRedisService(redisConfig models.RedisConfig) *RedisService {
	return &RedisService{
		config: redisConfig,
	}
}

func (s *RedisService) SaveRegSession(token, username, email, hashedPassword string, expireTime int64) error {
	regInfo := regInfo{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	if _, err := s.config.JsonHandler.JSONSet(token, ".", regInfo); err != nil {
		return err
	}

	if _, err := s.config.Client.Do(s.config.Ctx, "EXPIRE", token, expireTime).Result(); err != nil {
		return err
	}

	return nil
}

func (s *RedisService) GetRegSession(token string) (*regInfo, error) {
	jsonBytes, err := redis.Bytes(s.config.JsonHandler.JSONGet(token, "."))
	if err != nil {
		return nil, err
	}

	regInfo := regInfo{}
	if err = json.Unmarshal(jsonBytes, &regInfo); err != nil {
		return nil, err
	}

	return &regInfo, nil
}

func (s *RedisService) SavePassResetSession(token string, userId int64, expireTime time.Duration) error {
	if err := s.config.Client.Set(s.config.Ctx, token, strconv.Itoa(int(userId)), expireTime).Err(); err != nil {
		return err
	}
	return nil
}

func (s *RedisService) GetPassResetSession(token string) (int64, error) {
	userid, err := s.config.Client.Get(s.config.Ctx, token).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func (s *RedisService) SaveAuthSession(userId int64, td *models.TokenDetails) error {
	// converting Unix to UTC
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	if err := s.config.Client.Set(s.config.Ctx, td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err(); err != nil {
		return err
	}

	if err := s.config.Client.Set(s.config.Ctx, td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err(); err != nil {
		return err
	}

	return nil
}

func (s *RedisService) GetAuthSession(accessDetails *models.AccessDetails) (int64, error) {
	userid, err := s.config.Client.Get(s.config.Ctx, accessDetails.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func (s *RedisService) DeleteSession(key string) (int64, error) {
	deleted, err := s.config.Client.Del(s.config.Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
