package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trkhoi/khoi-api/internal/response"
)

type handler struct {
}

func New() IUserHandler {
	return &handler{}
}

func (h *handler) GetUserDetail(c *gin.Context) {
	resp := "Hello world"

	c.JSON(http.StatusOK, response.CreateResponse(resp, nil, nil, nil))
}
