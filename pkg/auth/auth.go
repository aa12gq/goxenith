package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goxenith/app/models/ent"
	"goxenith/pkg/logger"
)

// CurrentUser 从 gin.context 中获取当前登录用户
func CurrentUser(c *gin.Context) *ent.User {
	_user, ok := c.MustGet("current_user").(*ent.User)
	if !ok {
		logger.LogIf(errors.New("无法获取用户"))
		return &ent.User{}
	}
	return _user
}

// CurrentUID 从 gin.context 中获取当前登录用户 ID
func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}
