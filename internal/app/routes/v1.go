package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/trkhoi/khoi-api/internal/app/handler"
	"github.com/trkhoi/khoi-api/internal/appmain"
)

func LoadV1(r *gin.Engine, h *handler.Handler, p *appmain.Params) {
	v1 := r.Group("/api/v1")

	// v1.Use(middleware.WithTestCode())

	v1.GET("/user", h.User.GetUserDetail)
}
