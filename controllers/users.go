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
	var input models.UserCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf(input.Username)
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return
	}

	if res, err := DB.Exec("INSERT INTO `users` (`username`) VALUES (?)", input.Username); err != nil {
		driverErr, ok := err.(*mysql.MySQLError)
		if ok && driverErr.Number == 1062 {
			c.JSON(http.StatusUnprocessableEntity, models.Error{
				Success: false,
				Message: "User already exists with this username",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Error{
			Success: false,
			Message: "Error creating user",
		})
		log.Printf("Error creating user: %s\n", driverErr.Error())
	} else {
		id, _ := res.LastInsertId()
		c.JSON(http.StatusOK, models.UserResponse{
			Success: true,
			User: &models.User{
				ID:       id,
				Username: input.Username,
				LastSeen: time.Now().Unix(),
			},
		})
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
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnprocessableEntity, models.Error{
				Success: false,
				Message: "Invalid username",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, models.Error{
				Success: false,
				Message: "Error getting user",
			})
			log.Printf("Error getting user: %s\n", err.Error())
			return
		}
	}
	if rank.Valid {
		user.Rank = &rank.String
	}
	user.LastSeen = lastSeen.Time.Unix()

	// ban info
	var ban models.Ban
	var banTo sql.NullTime
	if err := DB.QueryRow("SELECT `ban_to`, `reason`, `admin` FROM `bans` WHERE `username` = ? AND `status` = true",
		username).Scan(&banTo, &ban.Reason, &ban.Admin); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, models.Error{
				Success: false,
				Message: "Error getting user's ban info",
			})
			log.Printf("Error getting user's ban info: %s\n", err.Error())
			return
		}
	} else {
		if banTo.Valid {
			banToInt64 := banTo.Time.Unix()
			ban.To = &banToInt64
		}
		user.Ban = &ban
	}

	// mute info
	var mute models.Ban
	var muteTo sql.NullTime
	if err := DB.QueryRow("SELECT `mute_to`, `reason`, `admin` FROM `mutes` WHERE `username` = ? AND `status` = true",
		username).Scan(&muteTo, &mute.Reason, &mute.Admin); err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, models.Error{
				Success: false,
				Message: "Error getting user's mute info",
			})
			log.Printf("Error getting user's mute info: %s\n", err.Error())
			return
		}
	} else {
		if muteTo.Valid {
			muteToInt64 := muteTo.Time.Unix()
			mute.To = &muteToInt64
		}
		user.Mute = &mute
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Success: true,
		User:    &user,
	})
}

func UpdateUserExp(c *gin.Context) {
	var input models.UserExpUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf(input.Username)
		c.JSON(http.StatusBadRequest, models.Error{
			Success: false,
			Message: "Missing one or more fields " + err.Error(),
		})
		return
	}

	if _, err := DB.Exec("UPDATE `users` SET `exp` = `exp` + ? WHERE `username` = ?", input.Exp, input.Username); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnprocessableEntity, models.Error{
				Success: false,
				Message: "Invalid username",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Error{
				Success: false,
				Message: "Error updating user exp",
			})
			log.Printf("Error updating user exp: %s\n", err.Error())
		}
	} else {
		c.JSON(http.StatusOK, models.UserResponse{
			Success: true,
		})
	}
}
