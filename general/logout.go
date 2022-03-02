package general

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func LogoutHandler(c *gin.Context) {
	// delete session from database
	claims := extractClaimsFromContext(c)
	userId := claims["id"]
	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	if _, err := DB.Exec("DELETE FROM `sessions` WHERE `userId` = ?  AND `expires_at` = ?", userId, exp); err != nil {
		log.Printf(ErrDeletingSessionInfo.Error()+": %s\n", err.Error())
		handleErr(c, http.StatusInternalServerError, ErrDeletingSessionInfo)
		return
	}

	// delete auth cookie
	if Middleware.SendCookie {
		if Middleware.CookieSameSite != 0 {
			c.SetSameSite(Middleware.CookieSameSite)
		}

		c.SetCookie(
			Middleware.CookieName,
			"",
			-1,
			"/",
			Middleware.CookieDomain,
			Middleware.SecureCookie,
			Middleware.CookieHTTPOnly,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
