package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
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

func (s *ErrService) HandleInternalErr(c *gin.Context, httpCode int, err, internalErr error) {
	log.Printf(err.Error()+": %s\n", internalErr.Error())
	s.HandleErr(c, httpCode, err)
}

func (s *ErrService) HandleDBErr(c *gin.Context, err error) {
	log.Printf(err.Error()+": %s\n", err.Error())
	s.HandleErr(c, http.StatusInternalServerError, errors.ErrDBQuery)
}
