package general

import (
	"context"
	"encoding/json"
	goredis "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	"strconv"
	"time"
)

var (
	redisClient      *goredis.Client
	redisJsonHandler *rejson.Handler
	ctx              context.Context
)

type regInfo struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
}

func SetupRedis(client *goredis.Client, jsonHandler *rejson.Handler) {
	redisClient = client
	redisJsonHandler = jsonHandler
	ctx = context.TODO()
}

func saveRegSession(token, username, email, hashedPassword string) error {
	regInfo := regInfo{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	if _, err := redisJsonHandler.JSONSet(token, ".", regInfo); err != nil {
		return err
	}

	if _, err := redisClient.Do(ctx, "EXPIRE", token, Middleware.RegTokenTime).Result(); err != nil {
		return err
	}

	return nil
}

func getRegSession(token string) (*regInfo, error) {
	jsonBytes, err := redis.Bytes(redisJsonHandler.JSONGet(token, "."))
	if err != nil {
		return nil, err
	}

	regInfo := regInfo{}
	if err = json.Unmarshal(jsonBytes, &regInfo); err != nil {
		return nil, err
	}

	return &regInfo, nil
}

func savePassResetSession(token string, userId int64) error {
	if err := redisClient.Set(ctx, token, strconv.Itoa(int(userId)), Middleware.PassResetTokenTime).Err(); err != nil {
		return err
	}
	return nil
}

func getPassResetSession(token string) (int64, error) {
	userid, err := redisClient.Get(ctx, token).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func saveAuthSession(userId int64, td *TokenDetails) error {
	// converting Unix to UTC
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	if err := redisClient.Set(ctx, td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err(); err != nil {
		return err
	}

	if err := redisClient.Set(ctx, td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err(); err != nil {
		return err
	}

	return nil
}

func getAuthSession(accessDetails *AccessDetails) (int64, error) {
	userid, err := redisClient.Get(ctx, accessDetails.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func deleteSession(key string) (int64, error) {
	deleted, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
