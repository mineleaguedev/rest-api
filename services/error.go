package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/models"
	"log"
	"net/http"
)

type ErrService struct {
}

func NewErrService() *ErrService {
	return &ErrService{}
}

func (s *ErrService) HandleErr(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, models.Error{
		Success: false,
		Message: err.Error(),
	})
}

func (s *ErrService) HandleInternalErr(c *gin.Context, err, internalErr error) {
	log.Printf(err.Error()+": %s\n", internalErr.Error())
	s.HandleErr(c, http.StatusInternalServerError, err)
}
