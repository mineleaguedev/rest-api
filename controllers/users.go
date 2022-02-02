package controllers

import (
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
