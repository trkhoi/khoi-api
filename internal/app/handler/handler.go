package handler

import (
	"github.com/trkhoi/khoi-api/internal/app/handler/user"
	"github.com/trkhoi/khoi-api/internal/appmain"
)

type Handler struct {
	User user.IUserHandler
}

func New(p *appmain.Params) *Handler {
	// svc := service.New(p)
	// controller := controller.New(p, svc)

	return &Handler{
		User: user.New(),
	}
}
