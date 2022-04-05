package paper

import (
	"github.com/mineleaguedev/rest-api/services"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{
		services: services,
	}
}