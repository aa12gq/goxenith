package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "goxenith/app/http/controllers/api/v1"
	"goxenith/app/models/ent"
	entUser "goxenith/app/models/ent/user"
	"goxenith/app/requests"
	"goxenith/dao"
	"goxenith/pkg/auth"
	"goxenith/pkg/logger"
	"goxenith/pkg/password"
	"goxenith/pkg/response"
	pb "goxenith/proto/app/v1"
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

	user, err := dao.DB.User.Query().Where(entUser.PhoneEQ(request.Phone)).First(c)

	if err != nil {
		response.Error(c, err, "账号不存在")
	} else {
		token := auth.NewJWT().IssueToken(user.ID, user.UserName)

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

	user, err := dao.DB.User.Query().Where(
		entUser.Or(
			entUser.PhoneEQ(request.Account),
			entUser.EmailEQ(request.Account),
			entUser.UserNameEQ(request.Account),
		),
	).First(c)
	if err != nil {
		if ent.IsNotFound(err) {
			response.Abort404(c, "账号不存在")
		}
		logger.Warn(fmt.Sprintf("未找到账号为 %v 的用户信息", request.Account))
		response.Abort500(c, "账号查询出错")
	}

	if !password.BcryptPasswordMatch(request.Password, user.Password) {
		response.Unauthorized(c, "账号不存在或密码错误")
	}
	token := auth.NewJWT().IssueToken(user.ID, user.UserName)
	response.JSON(c, pb.LoginByPhoneReply{
		Uid:   user.ID,
		Token: token,
	})
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {

	token, err := auth.NewJWT().RefreshToken(c)

	if err != nil {
		response.Error(c, err, "令牌刷新失败")
	} else {
		response.JSON(c, gin.H{
			"token": token,
		})
	}
}
