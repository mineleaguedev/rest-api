package general

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type RegTokenInfo struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
}

func saveRegTokenSession(token, username, email, hashedPassword string) error {
	regTokenInfo := RegTokenInfo{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	_, err := RedisJsonHandler.JSONSet(token, ".", regTokenInfo)
	if err != nil {
		return err
	}

	_, err = RedisClient.Do(ctx, "EXPIRE", token, Middleware.RegTokenTime).Result()
	if err != nil {
		return err
	}

	return nil
}

func getRegTokenSession(token string) (*RegTokenInfo, error) {
	jsonBytes, err := redis.Bytes(RedisJsonHandler.JSONGet(token, "."))
	if err != nil {
		return nil, err
	}

	regTokenInfo := RegTokenInfo{}
	err = json.Unmarshal(jsonBytes, &regTokenInfo)
	if err != nil {
		return nil, err
	}

	return &regTokenInfo, nil
}

func savePassResetTokenSession(token string, userId int64) error {
	// converting Unix to UTC
	expires := time.Now().Add(Middleware.PassResetTokenTime)
	now := time.Now()

	err := RedisClient.Set(ctx, token, strconv.Itoa(int(userId)), expires.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func getPassResetTokenSession(token string) (int64, error) {
	userid, err := RedisClient.Get(ctx, token).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func saveSession(userId int64, td *TokenDetails) error {
	// converting Unix to UTC
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := RedisClient.Set(ctx, td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := RedisClient.Set(ctx, td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func deleteSession(uuid string) (int64, error) {
	deleted, err := RedisClient.Del(ctx, uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func getSession(accessDetails *AccessDetails) (int64, error) {
	userid, err := RedisClient.Get(ctx, accessDetails.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}
