package general

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

func RefreshHandler(c *gin.Context) {
	userId, tokenString, expire, err := refreshToken(c)
	if err != nil {
		handleErr(c, http.StatusUnauthorized, err)
		return
	}

	refreshResponse(c, http.StatusOK, userId, tokenString, expire)
}

func checkIfTokenExpire(c *gin.Context) (jwt.MapClaims, error) {
	token, err := parseToken(c)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	iat := int64(claims["iat"].(float64))
	exp := int64(claims["exp"].(float64))

	if iat > time.Now().Add(-Middleware.AccessTokenTime).Unix() {
		return nil, ErrNotExpiredAccessToken
	}

	if exp > time.Now().Unix() {
		return nil, ErrExpiredRefreshToken
	}

	return claims, nil
}

func refreshToken(c *gin.Context) (float64, string, time.Time, error) {
	claims, err := checkIfTokenExpire(c)
	if err != nil {
		return 0, "", time.Now(), err
	}

	userId := claims["id"].(float64)
	oldExp := time.Unix(int64(claims["exp"].(float64)), 0)

	var ip string
	err = DB.QueryRow("SELECT INET_NTOA(`ip`) FROM `sessions` WHERE `userId` = ? AND `expires_at` = ?", userId, oldExp).Scan(&ip)
	if err != nil || ip == "" {
		return 0, "", time.Now(), ErrInvalidToken
	}

	if _, err := DB.Exec("DELETE FROM `sessions` WHERE `userId` = ?  AND `expires_at` = ?", userId, oldExp); err != nil {
		return 0, "", time.Now(), ErrDeletingSessionInfo
	}

	if ip != clientIp(c) {
		return 0, "", time.Now(), ErrChangedClientIp
	}

	// create token
	newToken, expire, err := createToken(claims)
	if err != nil {
		return 0, "", time.Now(), err
	}

	return userId, newToken, expire, nil
}

func refreshResponse(c *gin.Context, _ int, userId float64, token string, expire time.Time) {
	if _, err := DB.Exec("INSERT INTO `sessions` (`token`, `userId`, `ip`, `expires_at`) VALUES (?, ?, INET_ATON(?), ?)",
		token, userId, clientIp(c), expire); err != nil {
		log.Printf(ErrAddingSessionInfo.Error()+": %s\n", err.Error())
		handleErr(c, http.StatusInternalServerError, ErrAddingSessionInfo)
		return
	}

	// set cookie
	sendCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}
