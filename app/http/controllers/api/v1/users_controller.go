package v1

import (
	"github.com/gin-gonic/gin"
	"goxenith/pkg/auth"
	"goxenith/pkg/response"
)

type UsersController struct {
	BaseAPIController
}

// CurrentUser 当前登录用户信息
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	response.Data(c, userModel)
}
