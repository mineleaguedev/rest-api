package cabinet

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
)

func (h *Handler) TransferMoneyHandler(c *gin.Context) {
	var input models.TransferMoneyRequest

	if err := c.ShouldBind(&input); err != nil {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrMissingTransferMoneyValues)
		return
	}

	if response := h.services.VerifyCaptcha(input.Captcha); !response.Success {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrInvalidCaptcha)
		return
	}

	userId := c.GetInt64("userId")

	var money int64
	if err := h.db.QueryRow("SELECT `money` FROM `users` WHERE `id` = ?", userId).Scan(&money); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleDBErr(c, err)
		}
		return
	}

	if money < input.Amount {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrNotEnoughMoney)
		return
	}

	var toId int64
	if err := h.db.QueryRow("SELECT `id` FROM `users` WHERE `username` LIKE ?", input.Username).Scan(&toId); err != nil {
		if err == sql.ErrNoRows {
			h.services.HandleErr(c, http.StatusBadRequest, errors.ErrUserDoesNotExist)
		} else {
			h.services.HandleDBErr(c, err)
		}
		return
	}

	if _, err := h.db.Exec("UPDATE `users` SET `money` = `money` - ? WHERE `id` LIKE ?", input.Amount, userId); err != nil {
		h.services.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrMoneySubtraction, err)
		return
	}

	if _, err := h.db.Exec("UPDATE `users` SET `money` = `money` + ? WHERE `username` LIKE ?", input.Amount, input.Username); err != nil {
		h.services.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrMoneyAddition, err)
		return
	}

	if _, err := h.db.Exec("INSERT INTO `transfers` (`from_id`, `to_id`, `amount`) VALUES (?, ?, ?)",
		userId, toId, input.Amount); err != nil {
		h.services.HandleInternalErr(c, http.StatusInternalServerError, errors.ErrSavingTransferInfo, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
	})
}
