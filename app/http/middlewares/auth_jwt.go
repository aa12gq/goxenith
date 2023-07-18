package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	entUser "goxenith/app/models/ent/user"
	"goxenith/pkg/auth"
	"goxenith/pkg/config"
	"goxenith/pkg/database"
	"goxenith/pkg/model"
	"goxenith/pkg/response"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.NewJWT().ParserToken(c)
		if err != nil {
			response.Unauthorized(c, fmt.Sprintf("请查看 %v 相关的接口认证文档", config.GetString("app.name")))
			return
		}
		user, err := database.DB.User.Query().Where(entUser.IDEQ(claims.UserID), entUser.DeleteEQ(model.DeletedNo)).First(c)
		if err != nil {
			return
		}
		if user.ID == 0 {
			response.Unauthorized(c, "找不到对应用户，用户可能已删除")
			return
		}

		// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
		c.Set("current_user_id", user.ID)
		c.Set("current_user_name", user.UserName)
		c.Set("current_user", user)

		c.Next()
	}
}
