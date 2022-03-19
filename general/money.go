package general

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
)

func TransferMoneyHandler(c *gin.Context) {
	var input models.TransferMoneyRequest

	if err := c.ShouldBind(&input); err != nil {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrMissingTransferMoneyValues)
		return
	}

	if response := Service.VerifyCaptcha(input.Captcha); !response.Success {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var money int64
	if err := DB.QueryRow("SELECT `money` FROM `users` WHERE `id` = ?", userId).Scan(&money); err != nil {
		Service.HandleInternalErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist, err)
		return
	}

	if money < input.Amount {
		Service.HandleErr(c, http.StatusBadRequest, errors.ErrNotEnoughMoney)
		return
	}

	var toId int64
	if err := DB.QueryRow("SELECT `id` FROM `users` WHERE `username` LIKE ?", input.Username).Scan(&toId); err != nil {
		Service.HandleInternalErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist, err)
		return
	}

	if _, err := DB.Exec("UPDATE `users` SET `money` = `money` - ? WHERE `id` LIKE ?", input.Amount, userId); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrTransferringMoney, err)
		return
	}

	if _, err := DB.Exec("UPDATE `users` SET `money` = `money` + ? WHERE `username` LIKE ?", input.Amount, input.Username); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrTransferringMoney, err)
		return
	}

	if _, err := DB.Exec("INSERT INTO `transfers` (`from_id`, `to_id`, `amount`) VALUES (?, ?, ?)",
		userId, toId, input.Amount); err != nil {
		Service.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSavingTransferInfo, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}