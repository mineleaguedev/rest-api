package general

import (
	"strconv"
	"time"
)

func saveSession(userId int64, td *TokenDetails) error {
	// converting Unix to UTC
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := RedisClient.Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := RedisClient.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func deleteSession(uuid string) (int64, error) {
	deleted, err := RedisClient.Del(uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func getSession(accessDetails *AccessDetails) (int64, error) {
	userid, err := RedisClient.Get(accessDetails.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}
