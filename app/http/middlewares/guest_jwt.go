package middlewares

import (
	"github.com/gin-gonic/gin"
	"goxenith/pkg/auth"
	"goxenith/pkg/response"
)

// GuestJWT 强制使用游客身份访问
func GuestJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(c.GetHeader("Authorization")) > 0 {

			// 解析 token 成功，说明登录成功
			_, err := auth.NewJWT().ParserToken(c)
			if err == nil {
				response.Unauthorized(c, "请使用游客身份访问")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
