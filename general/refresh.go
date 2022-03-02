package general

import (
	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

func HelloHandler(c *gin.Context) {
	claims := ginJwt.ExtractClaims(c)
	user, _ := c.Get(JWTMiddleware.IdentityKey)
	c.JSON(200, gin.H{
		"id":       user.(*User).ID,
		"username": claims["username"],
	})
}

func RefreshHandler(c *gin.Context) {
	userId, tokenString, expire, err := refreshToken(c)
	if err != nil {
		unauthorized(c, http.StatusUnauthorized, JWTMiddleware.HTTPStatusMessageFunc(err, c))
		return
	}

	refreshResponse(c, http.StatusOK, userId, tokenString, expire)
}

func unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm="+JWTMiddleware.Realm)
	if !JWTMiddleware.DisabledAbort {
		c.Abort()
	}

	JWTMiddleware.Unauthorized(c, code, message)
}

func checkIfTokenExpire(c *gin.Context) (jwt.MapClaims, error) {
	token, err := JWTMiddleware.ParseToken(c)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))
	exp := int64(claims["exp"].(float64))

	if origIat > JWTMiddleware.TimeFunc().Add(-JWTMiddleware.Timeout).Unix() {
		return nil, ErrNotExpiredAccessToken
	}

	if exp > JWTMiddleware.TimeFunc().Unix() {
		return nil, ErrExpiredRefreshToken
	}

	return claims, nil
}

func clientIp(c *gin.Context) string {
	clientIp := c.ClientIP()
	if clientIp == "::1" {
		clientIp = "127.0.0.1"
	}
	return clientIp
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

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(JWTMiddleware.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := JWTMiddleware.TimeFunc().Add(JWTMiddleware.Timeout)
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = JWTMiddleware.TimeFunc().Unix()
	tokenString, err := newToken.SignedString(JWTMiddleware.Key)
	if err != nil {
		return 0, "", time.Now(), err
	}

	return userId, tokenString, expire, nil
}

func refreshResponse(c *gin.Context, _ int, userId float64, token string, expire time.Time) {
	if _, err := DB.Exec("INSERT INTO `sessions` (`token`, `userId`, `ip`, `expires_at`) VALUES (?, ?, INET_ATON(?), ?)",
		token, userId, clientIp(c), expire); err != nil {
		handleErr(c, http.StatusInternalServerError, ErrAddingSessionInfo)
		log.Printf(ErrAddingSessionInfo.Error()+": %s\n", err.Error())
		return
	}

	// set cookie
	if JWTMiddleware.SendCookie {
		expireCookie := JWTMiddleware.TimeFunc().Add(JWTMiddleware.CookieMaxAge)
		maxAge := int(expireCookie.Unix() - time.Now().Unix())

		if JWTMiddleware.CookieSameSite != 0 {
			c.SetSameSite(JWTMiddleware.CookieSameSite)
		}

		c.SetCookie(
			JWTMiddleware.CookieName,
			token,
			maxAge,
			"/",
			JWTMiddleware.CookieDomain,
			JWTMiddleware.SecureCookie,
			JWTMiddleware.CookieHTTPOnly,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}
