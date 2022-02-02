package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
	"time"
)

func CreateUser(c *gin.Context) {
	var input models.UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": models.Error{
				Code:    400,
				Message: "Invalid user data",
			},
		})
		return
	}

	if res, err := DB.Exec("INSERT INTO `users` (`username`) VALUES (?)", input.Username); err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": models.Error{
						Code:    400,
						Message: "User already exists",
					},
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": models.Error{
						Code:    500,
						Message: "Error creating user",
					},
				})
				log.Printf("Error creating user: %s\n", driverErr.Error())
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": models.Error{
					Code:    500,
					Message: "Error creating user",
				},
			})
			log.Printf("Error creating user: %s\n", err.Error())
		}
	} else {
		id, _ := res.LastInsertId()
		user := models.User{
			ID:       id,
			Username: input.Username,
			Rank:     "PLAYER",
			LastSeen: time.Now().Unix(),
		}
		c.JSON(http.StatusOK, user)
	}
}

func GetUser(c *gin.Context) {
	username := c.Param("name")

	var user = models.User{
		Username: username,
	}

	var rank sql.NullString
	var lastSeen sql.NullTime
	if err := DB.QueryRow("SELECT `id`, `exp`, `rank`, `playtime`, `last_seen` FROM `users` WHERE `username` = ?",
		username).Scan(&user.ID, &user.Exp, &rank, &user.Playtime, &lastSeen); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": models.Error{
					Code:    500,
					Message: "Error getting user",
				},
			})
			log.Printf("Error getting user: %s\n", err.Error())
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": models.Error{
					Code:    400,
					Message: "User does not exist",
				},
			})
			return
		}
	}

	if !rank.Valid {
		user.Rank = "PLAYER"
	} else {
		user.Rank = rank.String
	}
	user.LastSeen = lastSeen.Time.Unix()

	// ban info
	var ban models.Ban
	var banTo sql.NullTime
	if err := DB.QueryRow("SELECT `ban_to`, `reason`, `admin` FROM `bans` WHERE `username` = ? AND `status` = true", username).Scan(&banTo, &ban.Reason, &ban.Admin); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": models.Error{
					Code:    500,
					Message: "Error getting user's ban info",
				},
			})
			log.Printf("Error getting user's ban info: %s\n", err.Error())
			return
		}
	} else {
		if banTo.Valid {
			ban.To = banTo.Time.Unix()
		}
		user.Ban = &ban
	}

	// mute info
	var mute models.Ban
	var muteTo sql.NullTime
	if err := DB.QueryRow("SELECT `mute_to`, `reason`, `admin` FROM `mutes` WHERE `username` = ? AND `status` = true", username).Scan(&muteTo, &mute.Reason, &mute.Admin); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": models.Error{
					Code:    500,
					Message: "Error getting user's mute info",
				},
			})
			log.Printf("Error getting user's mute info: %s\n", err.Error())
			return
		}
	} else {
		if muteTo.Valid {
			mute.To = muteTo.Time.Unix()
		}
		user.Mute = &mute
	}

	c.JSON(http.StatusOK, user)
}
