package general

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
	"time"
)

func AuthHandler(c *gin.Context) {
	data, httpCode, err := authenticate(c)
	if err != nil {
		handleErr(c, httpCode, ErrAddingSessionInfo)
		return
	}

	// create token
	token, expire, err := createToken(getPayload(data))
	if err != nil {
		log.Printf(ErrFailedTokenCreation.Error()+": %s\n", err.Error())
		handleErr(c, http.StatusInternalServerError, ErrFailedTokenCreation)
		return
	}

	authResponse(c, http.StatusOK, token, expire)
}

func authenticate(c *gin.Context) (interface{}, int, error) {
	var input models.AuthRequest

	if err := c.ShouldBind(&input); err != nil {
		return nil, http.StatusBadRequest, ErrMissingAuthValues
	}

	if response := CaptchaClient.VerifyToken(input.Captcha); !response.Success {
		return nil, http.StatusBadRequest, ErrInvalidCaptcha
	}

	var id float64
	var hashedPassword string
	err := DB.QueryRow("SELECT `id`, `password_hash` FROM `users` WHERE `username` = ?", input.Username).Scan(&id, &hashedPassword)
	if err != nil && err == sql.ErrNoRows {
		return nil, http.StatusBadRequest, ErrUserDoesNotExist
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, hashedPassword)
	if err != nil {
		log.Printf(ErrUnhashingPassword.Error()+": %s\n", err.Error())
		return nil, http.StatusInternalServerError, ErrUnhashingPassword
	}

	if !match {
		return nil, http.StatusBadRequest, ErrWrongUsernameOrPassword
	}

	return &User{
		ID:       id,
		Username: input.Username,
	}, http.StatusOK, nil
}

func getPayload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			"id":       v.ID,
			"username": v.Username,
		}
	}
	return jwt.MapClaims{}
}

func extractClaims(tokenString string) (jwt.MapClaims, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return Middleware.Key, nil
	})
	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true, nil
	} else {
		return nil, false, nil
	}
}

func getCountOfSessions(userId float64) (int, int, error) {
	rows, err := DB.Query("SELECT COUNT(*) FROM `sessions` WHERE `userId` = ?", userId)
	if err != nil {
		log.Printf(ErrCountingUserSessions.Error()+": %s\n", err.Error())
		return 0, http.StatusInternalServerError, ErrCountingUserSessions
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Printf(ErrCountingUserSessions.Error()+": %s\n", err.Error())
			return 0, http.StatusInternalServerError, ErrCountingUserSessions
		}
	}

	return count, http.StatusOK, nil
}

func authResponse(c *gin.Context, _ int, token string, expire time.Time) {
	claims, ok, err := extractClaims(token)
	if !ok && err != nil {
		log.Printf(ErrExtractingClaims.Error()+": %s\n", err.Error())
		handleErr(c, http.StatusInternalServerError, ErrExtractingClaims)
		return
	}

	userId := claims["id"].(float64)

	sessionCount, httpCode, err := getCountOfSessions(userId)
	if err != nil {
		handleErr(c, httpCode, err)
		return
	}

	if sessionCount >= 5 {
		if _, err := DB.Exec("DELETE FROM `sessions` WHERE `userId` = ?", userId); err != nil {
			log.Printf(ErrDeletingSessionInfo.Error()+": %s\n", err.Error())
			handleErr(c, http.StatusInternalServerError, ErrDeletingSessionInfo)
			return
		}
	}

	ip := clientIp(c)
	if _, err := DB.Exec("INSERT INTO `sessions` (`token`, `userId`, `ip`, `expires_at`) VALUES (?, ?, INET_ATON(?), ?)",
		token, userId, ip, expire); err != nil {
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
