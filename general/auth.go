package general

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
	"time"
)

func Authenticate(c *gin.Context) (interface{}, error) {
	var input models.AuthRequest

	if err := c.ShouldBind(&input); err != nil {
		return nil, ErrMissingAuthValues
	}

	if response := CaptchaClient.VerifyToken(input.Captcha); !response.Success {
		return nil, ErrInvalidCaptcha
	}

	var id float64
	var hashedPassword string
	err := DB.QueryRow("SELECT `id`, `password_hash` FROM `users` WHERE `username` = ?", input.Username).Scan(&id, &hashedPassword)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrUserDoesNotExist
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, hashedPassword)
	if err != nil {
		return nil, ErrUnhashingPassword
	}

	if !match {
		return nil, ErrUnknown
	}

	return &User{
		ID:       id,
		Username: input.Username,
	}, nil
}

func Payload(data interface{}) ginJwt.MapClaims {
	if v, ok := data.(*User); ok {
		return ginJwt.MapClaims{
			"id":       v.ID,
			"username": v.Username,
		}
	}
	return ginJwt.MapClaims{}
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return JWTMiddleware.Key, nil
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

func LoginResponse(c *gin.Context, _ int, token string, expire time.Time) {
	claims, ok, err := extractClaims(token)

	if !ok && err != nil {
		handleErr(c, http.StatusInternalServerError, ErrExtractingClaims)
		log.Printf(ErrExtractingClaims.Error()+": %s\n", err.Error())
		return
	}

	userId := claims["id"].(float64)

	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	sessionCount, httpCode, err := getCountOfSessions(userId)
	if err != nil {
		handleErr(c, httpCode, err)
		return
	}

	if sessionCount >= 5 {
		if _, err := DB.Exec("DELETE FROM `sessions` WHERE `userId` = ?", userId); err != nil {
			handleErr(c, http.StatusInternalServerError, ErrDeletingSessionInfo)
			log.Printf(ErrDeletingSessionInfo.Error()+": %s\n", err.Error())
			return
		}
	}

	if _, err := DB.Exec("INSERT INTO `sessions` (`token`, `userId`, `ip`, `expires_at`) VALUES (?, ?, INET_ATON(?), ?)",
		token, userId, ip, expire); err != nil {
		handleErr(c, http.StatusInternalServerError, ErrAddingSessionInfo)
		log.Printf(ErrAddingSessionInfo.Error()+": %s\n", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}
