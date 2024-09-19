package user

import "github.com/gin-gonic/gin"

type IUserHandler interface {
	GetUserDetail(c *gin.Context)
}
