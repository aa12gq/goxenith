package auth

import (
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	entUser "goxenith/app/models/ent/user"
	"goxenith/app/requests"
	"goxenith/pkg/auth"
	"goxenith/pkg/database"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
	"strconv"
)

type LoginController struct {
	v1.BaseAPIController
}

// LoginByPhone 手机登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {

	request := pb.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	user, err := database.DB.User.Query().Where(entUser.PhoneEQ(request.Phone)).First(c)

	if err != nil {
		response.Error(c, err, "账号不存在")
	} else {
		token := auth.NewJWT().IssueToken(strconv.FormatUint(user.ID, 10), user.UserName)

		response.JSON(c, pb.LoginByPhoneReply{
			Uid:   user.ID,
			Token: token,
		})
	}
}

// LoginByPassword 多种方法登录，支持手机号、email 和用户名
func (lc *LoginController) LoginByPassword(c *gin.Context) {
	// 1. 验证表单
	request := pb.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
		return
	}

	user, err := database.DB.User.Query().Where(
		entUser.Or(
			entUser.PhoneEQ(request.Account),
			entUser.EmailEQ(request.Account),
			entUser.UserNameEQ(request.Account),
		),
	).First(c)

	if err != nil {
		response.Unauthorized(c, "账号不存在或密码错误")
	} else {
		token := auth.NewJWT().IssueToken(strconv.FormatUint(user.ID, 10), user.UserName)

		response.JSON(c, pb.LoginByPhoneReply{
			Uid:   user.ID,
			Token: token,
		})
	}
}
